package main

import(
	"errors"
	"net/http"

	"greenlight/internal/data"
	"greenlight/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request){
	
	var input struct{
		Name string `json:"name"`
		Email string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil{
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name: input.Name,
		Email: input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil{
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.
}
