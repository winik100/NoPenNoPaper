package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/winik100/NoPenNoPaper/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamicChain := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	mux.Handle("GET /{$}", dynamicChain.ThenFunc(app.home))
	mux.Handle("GET /signup", dynamicChain.ThenFunc(app.signup))
	mux.Handle("POST /signup", dynamicChain.ThenFunc(app.signupPost))
	mux.Handle("GET /login", dynamicChain.ThenFunc(app.login))
	mux.Handle("POST /login", dynamicChain.ThenFunc(app.loginPost))
	mux.Handle("POST /logout", dynamicChain.ThenFunc(app.logoutPost))

	protectedChain := dynamicChain.Append(app.requireAuthentication)
	mux.Handle("GET /create", protectedChain.ThenFunc(app.create))
	mux.Handle("POST /create", protectedChain.ThenFunc(app.createPost))

	standardChain := alice.New(app.recoverPanic, app.logRequest, headers)
	return standardChain.Then(mux)
}
