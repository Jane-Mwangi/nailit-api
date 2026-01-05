package main

import (
	"errors"
	"net/http"

	"github.com/Jane-Mwangi/nailit-api/internal/data"
	"github.com/Jane-Mwangi/nailit-api/internal/validator"
	"github.com/google/uuid"
)

func (app *application) createServiceTypesHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		ServiceID       uuid.UUID `json:"service_id"`
		Name            string    `json:"name"`
		Price           int       `json:"price"`
		DurationMinutes int       `json:"duration_minutes"`
		ImageURL        string    `json:"image_url"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	serviceType := &data.ServiceType{
		ServiceID:       input.ServiceID,
		Name:            input.Name,
		Price:           input.Price,
		DurationMinutes: input.DurationMinutes,
		ImageURL:        input.ImageURL,
	}

	v := validator.New()

	if data.ValidateServiceType(v, serviceType); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.ServiceTypes.Insert(serviceType)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateServiceType):

			app.failedValidationResponse(w, r, map[string]string{
				"name": "service type already exists",
			})
		default:

			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"service_type": serviceType}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) getServiceTypeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	service_type, err := app.models.ServiceTypes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"service_type": service_type}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
