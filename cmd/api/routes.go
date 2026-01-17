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
	router.HandlerFunc(http.MethodGet, "/v1/services", app.requirePermission("services:read", app.getAllServicesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/services", app.requirePermission("services:write", app.createServiceHandler))
	router.HandlerFunc(http.MethodGet, "/v1/services/:id", app.requirePermission("services:read", app.getServiceHandler))
	router.HandlerFunc(http.MethodPut, "/v1/services/:id", app.requirePermission("services:write", app.updateServiceHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/services/:id", app.requirePermission("services:write", app.deleteServiceHandler))

	// ServiceTypes
	router.HandlerFunc(http.MethodPost, "/v1/service-types", app.requirePermission("service-types:write", app.createServiceTypesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/service-types/:id", app.requirePermission("service-types:read", app.getServiceTypeHandler))
	router.HandlerFunc(http.MethodGet, "/v1/service-types", app.requirePermission("service-types:read", app.getAllServiceTypesHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/service-types/:id", app.requirePermission("service-types:write", app.updateServiceTypeHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/service-types/:id", app.requirePermission("service-types:write", app.deleteServiceTypeHandler))

	// Staff
	router.HandlerFunc(http.MethodPost, "/v1/staff", app.requirePermission("staff:write", app.createStaffHandler))
	router.HandlerFunc(http.MethodGet, "/v1/staff/:id", app.requirePermission("staff:read", app.getStaffByIdHandler))
	router.HandlerFunc(http.MethodGet, "/v1/staff", app.requirePermission("staff:read", app.getAllStaffHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/staff/:id", app.requirePermission("staff:write", app.updateStaffHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/staff/:id", app.requirePermission("staff:write", app.deleteStaffHandler))

	// users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodGet, "/v1/user/activate", app.testActivateHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}
