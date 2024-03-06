package main

import(
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "enviroment: %s", app.config.env)
	fmt.Fprintf(w, "versoin: %s", version)

}
