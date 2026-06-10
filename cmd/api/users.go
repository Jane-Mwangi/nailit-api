package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Jane-Mwangi/nailit-api/internal/data"
	"github.com/Jane-Mwangi/nailit-api/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:  input.Name,
		Email: input.Email,
		Role:  "customer",
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {

		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// generate token after user has been created in the DB
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// send background email and recover panics
	// app.background(func() {

	// 	data := map[string]interface{}{
	// 		"activationToken": token.Plaintext,
	// 		"userID":          user.ID,
	// 	}

	// 	err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
	// 	if err != nil {
	// 		app.logger.PrintError(err, nil)
	// 	}
	// })

	// err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// }

	app.background(func() {

		activationLink := fmt.Sprintf("%s/activate?token=%s", app.config.frontendURL, token.Plaintext)

		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
			"activationLink":  activationLink,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

}

// activating a user
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		TokenPlaintext string `json:"token"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// for testing
// func (app *application) testActivateHandler(w http.ResponseWriter, r *http.Request) {
// 	token := r.URL.Query().Get("token")
// 	if token == "" {
// 		http.Error(w, "missing token", http.StatusBadRequest)
// 		return
// 	}

// 	user, err := app.models.Users.GetForToken(data.ScopeActivation, token)
// 	if err != nil {
// 		http.Error(w, "invalid or expired token", http.StatusBadRequest)
// 		return
// 	}

// 	user.Activated = true
// 	err = app.models.Users.Update(user)
// 	if err != nil {
// 		http.Error(w, "could not activate user", http.StatusInternalServerError)
// 		return
// 	}

// 	// Delete all activation tokens for this user
// 	_ = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)

// 	fmt.Fprintf(w, "User %s activated successfully!", user.Email)
// }
