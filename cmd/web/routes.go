package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /create", app.create)
	mux.HandleFunc("POST /create", app.createPost)
	return mux
}
