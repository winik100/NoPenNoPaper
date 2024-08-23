package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/winik100/NoPenNoPaper/internal/core"
	"github.com/winik100/NoPenNoPaper/internal/models"
	"github.com/winik100/NoPenNoPaper/internal/validators"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.GetInt(r.Context(), authenticatedUserIdKey)
	userName := app.sessionManager.GetString(r.Context(), authenticatedUserNameKey)
	if userId == 0 || userName == "" {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	redirect := fmt.Sprintf("/users/%s", userName)
	http.Redirect(w, r, redirect, http.StatusPermanentRedirect)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userForm{}
	w.WriteHeader(http.StatusOK)
	app.render(w, r, "signup.tmpl.html", data)
}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {
	var form userForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "Name", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.NotBlank(form.Password), "Password", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.MinChars(form.Password, 8), "Password", "Passwort muss mindestens 8 Zeichen lang sein.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(w, r, "signup.tmpl.html", data)
		return
	}

	exists, err := app.users.Exists(form.Name)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if exists {
		form.CheckField(false, "Name", "Dieser Name ist bereits vergeben.")
		data := app.newTemplateData(r)
		data.Form = form
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(w, r, "signup.tmpl.html", data)
		return
	}

	_, err = app.users.Insert(form.Name, form.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Erfolgreich registriert! Bitte einloggen.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userForm{}
	w.WriteHeader(http.StatusOK)
	app.render(w, r, "login.tmpl.html", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	var form userForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "Name", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.NotBlank(form.Password), "Password", "Dieses Feld kann nicht leer sein.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(w, r, "login.tmpl.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Name, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddGenericError("Name und/oder Password sind falsch.")

			data := app.newTemplateData(r)
			data.Form = form
			w.WriteHeader(http.StatusUnprocessableEntity)
			app.render(w, r, "login.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), authenticatedUserIdKey, id)
	app.sessionManager.Put(r.Context(), authenticatedUserNameKey, form.Name)

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	redirect := fmt.Sprintf("/users/%s", form.Name)
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func (app *application) logoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), authenticatedUserIdKey)
	app.sessionManager.Remove(r.Context(), authenticatedUserNameKey)
	app.sessionManager.Remove(r.Context(), isAuthenticatedKey)
	app.sessionManager.Put(r.Context(), roleKey, core.RoleAnon)
	app.sessionManager.Put(r.Context(), "flash", "Erfolgreich ausgeloggt!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
