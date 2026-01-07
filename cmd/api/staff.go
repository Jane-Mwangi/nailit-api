package main

import (
	"errors"
	"net/http"

	"github.com/Jane-Mwangi/nailit-api/internal/data"
	"github.com/Jane-Mwangi/nailit-api/internal/validator"
)

func (app *application) createStaffHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		IsActive bool   `json:"is_active"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	staff := &data.Staff{
		Name:     input.Name,
		Email:    input.Email,
		IsActive: input.IsActive,
	}

	v := validator.New()

	if data.ValidateStaff(v, staff); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Staff.Insert(staff)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateStaff):

			app.failedValidationResponse(w, r, map[string]string{
				"email": "a staff member with this email already exists",
			})
		default:

			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"staff": staff}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
