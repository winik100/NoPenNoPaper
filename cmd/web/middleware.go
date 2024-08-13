package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/winik100/NoPenNoPaper/internal/models"
)

func headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		method := r.Method
		uri := r.URL.RequestURI()
		proto := r.Proto

		app.log.Info("received request", "ip", ip, "method", method, "uri", uri, "version", proto)
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})
	return csrfHandler
}

const isAuthenticatedKey = "isAuthenticated"
const authenticatedUserIdKey = "authenticatedUserID"
const authenticatedUserNameKey = "authenticatedUserName"
const characterIdKey = "characterId"
const roleKey = "role"
const roleAnon = "anonymous"
const rolePlayer = "player"
const roleGM = "gm"

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userName := app.sessionManager.GetString(r.Context(), authenticatedUserNameKey)
		if userName == "" {
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.users.Get(userName)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				next.ServeHTTP(w, r)
				return
			}
			app.serverError(w, r, err)
			return
		}

		app.sessionManager.Put(r.Context(), isAuthenticatedKey, true)
		app.sessionManager.Put(r.Context(), roleKey, user.Role)

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		//dont allow caching for pages requiring authentication
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthorized(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
