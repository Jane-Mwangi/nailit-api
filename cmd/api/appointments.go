package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Jane-Mwangi/nailit-api/internal/data"
	"github.com/Jane-Mwangi/nailit-api/internal/validator"
	"github.com/google/uuid"
)

func (app *application) createAppointmentHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		ServiceTypeID uuid.UUID              `json:"service_type_id"`
		StartsAt      time.Time              `json:"starts_at"`
		EndsAt        time.Time              `json:"ends_at"`
		Status        data.AppointmentStatus `json:"status"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	serviceType, err := app.models.ServiceTypes.Get(input.ServiceTypeID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	appointment := &data.Appointment{
		CustomerID:    user.ID,
		ServiceTypeID: input.ServiceTypeID,
		ServiceID:     serviceType.ServiceID,
		StartsAt:      input.StartsAt.UTC(),
		EndsAt:        input.EndsAt.UTC(),
		Status:        input.Status,
	}

	v := validator.New()
	data.ValidateAppointment(v, appointment)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Appointments.Insert(appointment)
	if err != nil {
		switch {

		case errors.Is(err, data.ErrOverlappingAppointment):
			app.conflictResponse(w, r, errors.New("you already have an appointment at this time"))

		case errors.Is(err, data.ErrDuplicateAppointment):
			app.conflictResponse(w, r, errors.New("appointment already exists"))

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{
		"appointment": appointment,
	}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	appointment, err := app.models.Appointments.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"appointment": appointment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAllAppointmentsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "starts_at")
	input.Filters.SortSafelist = []string{
		"starts_at",
		"ends_at",
		"status",
		"created_at",
		"id",
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	appointments, metadata, err := app.models.Appointments.GetAll(input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{
		"appointments": appointments,
		"metadata":     metadata,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	appointment, err := app.models.Appointments.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		ServiceTypeID *uuid.UUID              `json:"service_type_id"`
		StartsAt      *time.Time              `json:"starts_at"`
		EndsAt        *time.Time              `json:"ends_at"`
		Status        *data.AppointmentStatus `json:"status"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.ServiceTypeID != nil {
		appointment.ServiceTypeID = *input.ServiceTypeID
	}
	if input.StartsAt != nil {
		appointment.StartsAt = *input.StartsAt
	}
	if input.EndsAt != nil {
		appointment.EndsAt = *input.EndsAt
	}
	if input.Status != nil {
		appointment.Status = *input.Status
	}

	v := validator.New()
	if data.ValidateAppointment(v, appointment); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Appointments.Update(appointment)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		case errors.Is(err, data.ErrOverlappingAppointment):
			app.conflictResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"appointment": appointment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Appointments.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "appointment deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
