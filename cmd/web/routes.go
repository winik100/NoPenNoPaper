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

	protectedChain := dynamicChain.Append(app.requireAuthentication)
	mux.Handle("GET /create", protectedChain.ThenFunc(app.create))
	mux.Handle("POST /create", protectedChain.ThenFunc(app.createPost))
	mux.Handle("GET /characters/{id}/delete", protectedChain.ThenFunc(app.delete))
	mux.Handle("POST /characters/{id}/delete", protectedChain.ThenFunc(app.deletePost))
	mux.Handle("GET /characters/{id}", protectedChain.ThenFunc(app.viewCharacter))
	mux.Handle("POST /characters/{id}/editStat", protectedChain.ThenFunc(app.editStat))
	mux.Handle("GET /characters/{id}/addSkill", protectedChain.ThenFunc(app.addSkill))
	mux.Handle("POST /characters/{id}/addSkill", protectedChain.ThenFunc(app.addSkillPost))
	mux.Handle("GET /characters/{id}/editSkill", protectedChain.ThenFunc(app.editSkill))
	mux.Handle("POST /characters/{id}/editSkill", protectedChain.ThenFunc(app.editSkillPost))
	mux.Handle("GET /characters/{id}/addCustomSkill", protectedChain.ThenFunc(app.addCustomSkill))
	mux.Handle("POST /characters/{id}/addCustomSkill", protectedChain.ThenFunc(app.addCustomSkillPost))
	mux.Handle("GET /characters/{id}/editCustomSkill", protectedChain.ThenFunc(app.editCustomSkill))
	mux.Handle("POST /characters/{id}/editCustomSkill", protectedChain.ThenFunc(app.editCustomSkillPost))
	mux.Handle("GET /characters/{id}/addItem", protectedChain.ThenFunc(app.addItem))
	mux.Handle("POST /characters/{id}/addItem", protectedChain.ThenFunc(app.addItemPost))
	mux.Handle("POST /characters/{id}/editItemCount", protectedChain.ThenFunc(app.editItemCount))
	mux.Handle("POST /characters/{id}/deleteItem", protectedChain.ThenFunc(app.deleteItemPost))
	mux.Handle("GET /characters/{id}/addNote", protectedChain.ThenFunc(app.addNote))
	mux.Handle("POST /characters/{id}/addNote", protectedChain.ThenFunc(app.addNotePost))
	mux.Handle("POST /characters/{id}/deleteNote", protectedChain.ThenFunc(app.deleteNotePost))
	mux.Handle("GET /users/{id}/uploadMaterial", protectedChain.ThenFunc(app.uploadMaterial))
	mux.Handle("POST /users/{id}/uploadMaterial", protectedChain.ThenFunc(app.uploadMaterialPost))

	//some helpers
	mux.Handle("GET /customSkillInput", protectedChain.ThenFunc(app.customSkillInput))

	standardChain := alice.New(app.recoverPanic, app.logRequest, headers)
	return standardChain.Then(mux)
}

func (app *application) routesNoMW() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /signup", app.signup)
	mux.HandleFunc("POST /signup", app.signupPost)
	mux.HandleFunc("GET /login", app.login)
	mux.HandleFunc("POST /login", app.loginPost)
	mux.HandleFunc("POST /logout", app.logoutPost)

	mux.HandleFunc("GET /create", app.create)
	mux.HandleFunc("POST /create", app.createPost)
	mux.HandleFunc("GET /characters/{id}/delete", app.delete)
	mux.HandleFunc("POST /characters/{id}/delete", app.deletePost)
	mux.HandleFunc("GET /characters/{id}", app.viewCharacter)
	mux.HandleFunc("GET /characters/{id}/editStat", app.editStat)
	mux.HandleFunc("GET /characters/{id}/addSkill", app.addSkill)
	mux.HandleFunc("POST /characters/{id}/addSkill", app.addSkillPost)
	mux.HandleFunc("GET /characters/{id}/editSkill", app.editSkill)
	mux.HandleFunc("POST /characters/{id}/editSkill", app.editSkillPost)
	mux.HandleFunc("GET /characters/{id}/addCustomSkill", app.addCustomSkill)
	mux.HandleFunc("POST /characters/{id}/addCustomSkill", app.addCustomSkillPost)
	mux.HandleFunc("GET /characters/{id}/editCustomSkill", app.editCustomSkill)
	mux.HandleFunc("POST /characters/{id}/editCustomSkill", app.editCustomSkillPost)
	mux.HandleFunc("GET /characters/{id}/addItem", app.addItem)
	mux.HandleFunc("POST /characters/{id}/addItem", app.addItemPost)
	mux.HandleFunc("POST /characters/{id}/deleteItem", app.deleteItemPost)
	mux.HandleFunc("GET /characters/{id}/addNote", app.addNote)
	mux.HandleFunc("POST /characters/{id}/addNote", app.addNotePost)
	mux.HandleFunc("POST /characters/{id}/deleteNote", app.deleteNotePost)

	//some helpers
	mux.HandleFunc("GET /customSkillInput", app.customSkillInput)

	standardChain := alice.New(app.recoverPanic, app.logRequest, headers)
	return standardChain.Then(mux)
}
