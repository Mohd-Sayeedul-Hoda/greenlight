package main

import(
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router{
	
	router := httprouter.New()

	// adding custom error handler for notfound
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// adding custom error handler for method not allow
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)


	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)

	return router
}
