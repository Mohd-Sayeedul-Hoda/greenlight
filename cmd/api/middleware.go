package main

import (
	"errors"
	"strings"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"greenlight/internal/data"
	"greenlight/internal/validator"
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

	type client struct{
		limiter *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu sync.Mutex
		clients = make(map[string]*client)
	)

	go func(){
		for{
			// we are runnig this func every min
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients{
				if time.Since(client.lastSeen) > 3*time.Minute{
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
		}()

	// closure function 
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		if app.config.limiter.enabled{

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil{
				app.serverErrorResponse(w, r, err)
				return
			}

			mu.Lock()

			if _, found := clients[ip]; !found{
				// 2 is num per second at which bucket will fill and 4 is total no req allow
				clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow(){
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}
			mu.Unlock()
		}

		next.ServeHTTP(w, r)
		})
}

func (app *application) authenticate(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// Add the "vary authorization header to the response"
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		// if no authorization found we make anonymous user add to the request
		if authorizationHeader == ""{
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer"{
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]
		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid(){
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil{
			switch{
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)

		next.ServeHTTP(w, r)
		})
}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc{

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// information from the request context
		user := app.contextGetUser(r)

		// If user is anonymous then call the authenticationRequired
		if user.IsAnonymous(){
			app.authenticationRequiredRespones(w, r)
			return
		}
		next.ServeHTTP(w, r)
		})
}

// so we 
func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc{
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		user := app.contextGetUser(r)

		if !user.Activated{
			app.inactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
		})

	return app.requireAuthenticatedUser(fn)
}
func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc{
	fn := func(w http.ResponseWriter, r *http.Request){
		// Retrieve the user from the request context
		user := app.contextGetUser(r)

		permissions, err := app.models.Permissions.GetAllForUser(user.ID)
		if err != nil{
			app.serverErrorResponse(w, r, err)
			return
		}

		if !permissions.Include(code){
			app.notPermittedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
	return app.requireActivatedUser(fn)
}
