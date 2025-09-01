package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	//setup up a new router 
	router := httprouter.New()
	//Handle 404
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	//Handle 405
	router.MethodNotAllowed = http.HandlerFunc(app.methosNotAllowedResponse)
	//setup routes
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)


	return app.recoverPanic(router)

}