package main

import (
	"errors"
	"fmt"
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

func (app *application) getStaffByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	staff, err := app.models.Staff.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"staff": staff}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) updateStaffHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	staff, err := app.models.Staff.Get(id)
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
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		IsActive *bool   `json:"is_active"`
		Version  int     `json:"version"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name == nil &&
		input.Email == nil &&
		input.IsActive == nil {
		app.badRequestResponse(w, r, errors.New("body must contain at least one updatable field"))
		return
	}

	if input.Name != nil {
		staff.Name = *input.Name
	}
	if input.Email != nil {
		staff.Email = *input.Email
	}
	if input.IsActive != nil {
		staff.IsActive = *input.IsActive
	}

	v := validator.New()

	if data.ValidateStaff(v, staff); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Staff.Update(staff)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"staff": staff}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteStaffHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Staff.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "staff successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAllStaffHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string
		Email    string
		IsActive bool
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Email = app.readString(qs, "email", "")
	input.IsActive = app.readBool(qs, "is_active", true)

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "name", "email", "is_active", "-name", "-email", "-is_active"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	staff, metadata, err := app.models.Staff.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"staff": staff, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}
