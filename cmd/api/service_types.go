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

	err = app.writeJSON(w, http.StatusCreated, envelope{"serviceType": serviceType}, nil)
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

	serviceType, err := app.models.ServiceTypes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"serviceType": serviceType}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) updateServiceTypeHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	serviceType, err := app.models.ServiceTypes.Get(id)
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
		Name            *string `json:"name"`
		Price           *int    `json:"price"`
		DurationMinutes *int    `json:"duration_minutes"`
		ImageURL        *string `json:"image_url"`
		Version         int     `json:"version"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Reject empty PATCH
	if input.Name == nil &&
		input.Price == nil &&
		input.DurationMinutes == nil &&
		input.ImageURL == nil {
		app.badRequestResponse(w, r, errors.New("body must contain at least one updatable field"))
		return
	}

	if input.Name != nil {
		serviceType.Name = *input.Name
	}
	if input.Price != nil {
		serviceType.Price = *input.Price
	}
	if input.DurationMinutes != nil {
		serviceType.DurationMinutes = *input.DurationMinutes
	}
	if input.ImageURL != nil {
		serviceType.ImageURL = *input.ImageURL
	}

	serviceType.Version = input.Version

	v := validator.New()

	if data.ValidateServiceType(v, serviceType); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.ServiceTypes.Update(serviceType)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"serviceType": serviceType}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteServiceTypeHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.ServiceTypes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "service-type successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAllServiceTypesHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name            string
		Price           int
		DurationMinutes int
		ImageURL        string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Price = app.readInt(qs, "price", 0, v)
	input.DurationMinutes = app.readInt(qs, "duration_minutes", 0, v)
	input.ImageURL = app.readString(qs, "image_url", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "name", "price", "duration_minutes", "-name", "-price", "-duration_minutes"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	serviceTypes, metadata, err := app.models.ServiceTypes.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"serviceTypes": serviceTypes, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}


}
