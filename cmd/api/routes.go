package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// Services
	router.HandlerFunc(http.MethodGet, "/v1/services", app.requireActivatedUser(app.getAllServicesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/services", app.requireActivatedUser(app.createServiceHandler))
	router.HandlerFunc(http.MethodGet, "/v1/services/:id", app.requireActivatedUser(app.getServiceHandler))
	router.HandlerFunc(http.MethodPut, "/v1/services/:id", app.requireActivatedUser(app.updateServiceHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/services/:id", app.requireActivatedUser(app.deleteServiceHandler))

	// ServiceTypes
	router.HandlerFunc(http.MethodPost, "/v1/service-types", app.requireActivatedUser(app.createServiceTypesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/service-types/:id", app.requireActivatedUser(app.getServiceTypeHandler))
	router.HandlerFunc(http.MethodGet, "/v1/service-types", app.requireActivatedUser(app.getAllServiceTypesHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/service-types/:id", app.requireActivatedUser(app.updateServiceTypeHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/service-types/:id", app.requireActivatedUser(app.deleteServiceTypeHandler))

	// Staff
	router.HandlerFunc(http.MethodPost, "/v1/staff", app.requireActivatedUser(app.createStaffHandler))
	router.HandlerFunc(http.MethodGet, "/v1/staff/:id", app.requireActivatedUser(app.getStaffByIdHandler))
	router.HandlerFunc(http.MethodGet, "/v1/staff", app.requireActivatedUser(app.getAllStaffHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/staff/:id", app.requireActivatedUser(app.updateStaffHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/staff/:id", app.requireActivatedUser(app.deleteStaffHandler))

	// users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodGet, "/v1/user/activate", app.testActivateHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
