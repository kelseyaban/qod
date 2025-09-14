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
	router.HandlerFunc(http.MethodPost, "/v1/quotes", app.createQuoteHandler)
	router.HandlerFunc(http.MethodGet, "/v1/quotes/:id", app.displayQuoteHandler)
	router.HandlerFunc(http.MethodPatch,"/v1/quotes/:id",app.updateQuoteHandler)


	return app.recoverPanic(router)

}