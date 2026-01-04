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
		case errors.Is(err, data.ErrDuplicateService):

			app.failedValidationResponse(w, r, map[string]string{
				"name": "a service with this name already exists",
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
