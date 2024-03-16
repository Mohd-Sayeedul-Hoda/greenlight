package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	
		defer func(){
			// builtin recover function to check if there has been a panic or not
			if err := recover(); err != nil{
				// this will make go http server automatically close
				w.Header().Set("Connection", "close")

				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler)http.Handler{
	// this is total rate limit means golbally valid for all user
	// 2 means token bucket fill with 2 request per second
	// with a maximum of 4 request in a single burst
	limiter := rate.NewLimiter(2, 4)

	// clousre function 
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if !limiter.Allow(){
			app.rateLimitExceededResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
