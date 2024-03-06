package main

import(
	"fmt"
	"net/http"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "create movies\n")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request){
	
	id, err := app.readIDParam(r)
	if err != nil{
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "show the detail of movies %d\n", id)
}

