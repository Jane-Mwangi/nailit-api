package main

import (
	"net/http"
	"time"

	"github.com/Jane-Mwangi/nailit-api/internal/data"
	"github.com/Jane-Mwangi/nailit-api/internal/validator"
)

func (app *application) createServiceHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name string `json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	service := &data.Service{
		Name: input.Name,
	}

	v := validator.New()

	if data.ValidateService(v, service); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Services.Insert(service)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"service": service}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getServiceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	service := data.Service{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      "Casablanca",
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"service": service}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
