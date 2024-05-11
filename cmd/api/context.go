package main

import(
	"context"
	"net/http"

	"greenlight/internal/data"
)

type contextKey string

//Convert the string "user" to a contextKey type and assign it to user Contextkey
const userContextKey = contextKey("user")


// we add user struct to request as context
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request{
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User{
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}



