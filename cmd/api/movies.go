package main

import(
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "create movies\n")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request){

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1{
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "show the detail of movies %d\n", id)
}

