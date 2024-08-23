package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/winik100/NoPenNoPaper/internal/core"
	"github.com/winik100/NoPenNoPaper/internal/models"
	"github.com/winik100/NoPenNoPaper/internal/validators"
)

type userForm struct {
	Name                     string
	Password                 string
	validators.FormValidator `schema:"-"`
}

type uploadForm struct {
	FileName                 string
	Title                    string
	UploadedById             int
	UploadedByName           string
	validators.FormValidator `schema:"-"`
}

func (app *application) user(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.GetInt(r.Context(), authenticatedUserIdKey)
	userName := app.sessionManager.GetString(r.Context(), authenticatedUserNameKey)

	user, err := app.users.Get(userName)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.User = user
	if user.IsGM() {
		characters, err := app.characters.GetAll()
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		data.Characters = characters
	} else {
		characters, err := app.characters.GetAllFrom(userId)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		data.Characters = characters
	}

	w.WriteHeader(http.StatusOK)
	app.render(w, r, "user.tmpl.html", data)
}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
	userName := app.sessionManager.GetString(r.Context(), authenticatedUserNameKey)

	tmplStr := `<form id="deleteUserForm" action="/users/{{.Form.Name}}/delete" method="POST">
					<p id="deleteUserMessage">Sicher? Kann nicht rückgängig gemacht werden!</p>
					<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
					<input type="hidden" name="Name" Value="{{.Form.Name}}">
					<button type="submit">OK</button>
					<button hx-get="/users/{{.Form.Name}}" hx-target="#deleteUserForm" hx-select="#deleteUser" hx-swap="outerHTML">Abbrechen</button>
            	</form>`

	data := app.newTemplateData(r)
	data.Form = map[string]any{
		"Name": userName,
	}
	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "deleteUser", tmplStr, data)
}

func (app *application) deleteUserPost(w http.ResponseWriter, r *http.Request) {
	type deleteForm struct {
		Name string
	}

	var form deleteForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	err = app.users.Delete(form.Name)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), authenticatedUserIdKey)
	app.sessionManager.Remove(r.Context(), authenticatedUserNameKey)
	app.sessionManager.Remove(r.Context(), isAuthenticatedKey)
	app.sessionManager.Put(r.Context(), roleKey, core.RoleAnon)
	app.sessionManager.Put(r.Context(), "flash", "Nutzerkonto erfolgreich gelöscht!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) uploadMaterial(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.GetInt(r.Context(), authenticatedUserIdKey)
	userName := app.sessionManager.GetString(r.Context(), authenticatedUserNameKey)

	data := app.newTemplateData(r)
	data.Form = uploadForm{
		UploadedById:   userId,
		UploadedByName: userName,
	}
	w.WriteHeader(http.StatusOK)
	app.render(w, r, "upload.tmpl.html", data)
}

func (app *application) uploadMaterialPost(w http.ResponseWriter, r *http.Request) {
	var form uploadForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	file, header, err := r.FormFile("FilePath")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer file.Close()

	err = app.users.AddMaterial(form.Title, header.Filename, form.UploadedById)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateFileName) {
			form.AddGenericError("Eine Datei mit diesem Namen existiert bereits.")

			data := app.newTemplateData(r)
			data.Form = form
			w.WriteHeader(http.StatusUnprocessableEntity)
			app.render(w, r, "upload.tmpl.html", data)
			return
		}
		app.serverError(w, r, err)
		return
	}

	path := filepath.Join("ui/static/img/uploads/", strconv.Itoa(form.UploadedById))
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	f, err := os.Create(filepath.Join(path, header.Filename))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	redirect := fmt.Sprintf("/users/%s", form.UploadedByName)
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func (app *application) deleteMaterial(w http.ResponseWriter, r *http.Request) {
	var form uploadForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	err = app.users.DeleteMaterial(form.FileName, form.UploadedById)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	path := filepath.Join("ui/static/img/uploads/", strconv.Itoa(form.UploadedById), form.FileName)
	err = os.Remove(path)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	app.renderHtmx(w, r, "deleteMaterial", "", data)
}
