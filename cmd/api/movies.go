package main

import(
	"fmt"
	"time"
	"net/http"

	"greenlight/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "create movies\n")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request){
	
	id, err := app.readIDParam(r)
	if err != nil{
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID: id,
		CreatedAt: time.Now(),
		Title: "Casablaca", 
		Runtime: 102,
		Genres: []string{"drama", "romance", "war"},
		Version: 1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil{
		app.logger.Print(err)
		app.serverErrorResponse(w, r, err)
	}
}

