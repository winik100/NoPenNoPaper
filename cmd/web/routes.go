package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/winik100/NoPenNoPaper/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamicChain := alice.New(app.sessionManager.LoadAndSave, app.authenticate, noSurf)

	mux.Handle("GET /{$}", dynamicChain.ThenFunc(app.home))
	mux.Handle("GET /signup", dynamicChain.ThenFunc(app.signup))
	mux.Handle("POST /signup", dynamicChain.ThenFunc(app.signupPost))
	mux.Handle("GET /login", dynamicChain.ThenFunc(app.login))
	mux.Handle("POST /login", dynamicChain.ThenFunc(app.loginPost))
	mux.Handle("POST /logout", dynamicChain.ThenFunc(app.logoutPost))

	protectedChain := dynamicChain.Append(app.requireAuthentication, app.restrict)
	mux.Handle("GET /create", protectedChain.ThenFunc(app.create))
	mux.Handle("POST /create", protectedChain.ThenFunc(app.createPost))
	mux.Handle("GET /characters/{id}", protectedChain.ThenFunc(app.viewCharacter))
	mux.Handle("GET /characters/{id}/addItem", protectedChain.ThenFunc(app.addItem))
	mux.Handle("POST /characters/{id}/addItem", protectedChain.ThenFunc(app.addItemPost))
	mux.Handle("GET /characters/{id}/addNote", protectedChain.ThenFunc(app.addNote))
	mux.Handle("POST /characters/{id}/addNote", protectedChain.ThenFunc(app.addNotePost))

	//some helpers
	mux.Handle("POST /inc", protectedChain.ThenFunc(app.Inc))
	mux.Handle("POST /dec", protectedChain.ThenFunc(app.Dec))
	mux.Handle("GET /customSkillInput", protectedChain.ThenFunc(app.customSkillInput))
	mux.Handle("GET /cancel", protectedChain.ThenFunc(app.cancel))

	standardChain := alice.New(app.recoverPanic, app.logRequest, headers)
	return standardChain.Then(mux)
}

const (
	RoleAnon   string = "anonymous"
	RolePlayer string = "player"
	RoleGM     string = "gm"
)

var Permissions = map[string][]string{
	RoleAnon:   {"/", "/signup", "/login"},
	RolePlayer: {"/", "/signup", "/login", "/logout", "/create", "/characters/.*", "/characters/.*/addItem", "/characters/.*/addNote", "/inc", "/dec", "/customSkillInput"},
	RoleGM:     {"/", "/signup", "/login", "/logout", "/create", "/characters/.*", "/characters/.*/addItem", "/characters/.*/addNote", "/inc", "/dec", "/customSkillInput"},
}
