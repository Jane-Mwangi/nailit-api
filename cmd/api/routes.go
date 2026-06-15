package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/services", app.requirePermission("services:read", app.getAllServicesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/services", app.requirePermission("services:write", app.createServiceHandler))
	router.HandlerFunc(http.MethodGet, "/v1/services/:id", app.requirePermission("services:read", app.getServiceHandler))
	router.HandlerFunc(http.MethodPut, "/v1/services/:id", app.requirePermission("services:write", app.updateServiceHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/services/:id", app.requirePermission("services:write", app.deleteServiceHandler))

	router.HandlerFunc(http.MethodPost, "/v1/service-types", app.requirePermission("service-types:write", app.createServiceTypesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/service-types/:id", app.requirePermission("service-types:read", app.getServiceTypeHandler))
	router.HandlerFunc(http.MethodGet, "/v1/service-types", app.requirePermission("service-types:read", app.getAllServiceTypesHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/service-types/:id", app.requirePermission("service-types:write", app.updateServiceTypeHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/service-types/:id", app.requirePermission("service-types:write", app.deleteServiceTypeHandler))

	router.HandlerFunc(http.MethodPost, "/v1/staff", app.requirePermission("staff:write", app.createStaffHandler))
	router.HandlerFunc(http.MethodGet, "/v1/staff/:id", app.requirePermission("staff:read", app.getStaffByIdHandler))
	router.HandlerFunc(http.MethodGet, "/v1/staff", app.requirePermission("staff:read", app.getAllStaffHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/staff/:id", app.requirePermission("staff:write", app.updateStaffHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/staff/:id", app.requirePermission("staff:write", app.deleteStaffHandler))

	//appointments
	router.HandlerFunc(http.MethodPost, "/v1/appointments", app.requirePermission("appointments:write", app.createAppointmentHandler))
	router.HandlerFunc(http.MethodGet, "/v1/appointments/:id", app.requirePermission("appointments:read", app.getAppointmentHandler))
	router.HandlerFunc(http.MethodGet, "/v1/appointments", app.requirePermission("appointments:read", app.getAllAppointmentsHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/appointments/:id", app.requirePermission("appointments:write", app.updateAppointmentHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/appointments/:id", app.requirePermission("appointments:write", app.deleteAppointmentHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.Handler(http.MethodGet, "/metrics", promhttp.Handler())

	// metrics wraps ratelimit and auth as I want observe rejected requests too
	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(router))))
}
