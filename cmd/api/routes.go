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
	router.HandlerFunc(http.MethodGet, "/v1/services", app.getAllServicesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/services", app.createServiceHandler)
	router.HandlerFunc(http.MethodGet, "/v1/services/:id", app.getServiceHandler)
	router.HandlerFunc(http.MethodPut, "/v1/services/:id", app.updateServiceHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/services/:id", app.deleteServiceHandler)

	// ServiceTypes
	router.HandlerFunc(http.MethodPost, "/v1/service-types", app.createServiceTypesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/service-types/:id", app.getServiceTypeHandler)
	router.HandlerFunc(http.MethodGet, "/v1/service-types", app.getAllServiceTypesHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/service-types/:id", app.updateServiceTypeHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/service-types/:id", app.deleteServiceTypeHandler)

	// Staff
	router.HandlerFunc(http.MethodPost, "/v1/staff", app.createStaffHandler)
	router.HandlerFunc(http.MethodGet, "/v1/staff/:id", app.getStaffByIdHandler)
	router.HandlerFunc(http.MethodGet, "/v1/staff", app.getAllStaffHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/staff/:id", app.updateStaffHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/staff/:id", app.deleteStaffHandler)

	return app.recoverPanic(router)
}
