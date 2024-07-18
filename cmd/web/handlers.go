package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/winik100/NoPenNoPaper/internal/models"
	"github.com/winik100/NoPenNoPaper/internal/validators"
)

type characterCreateForm struct {
	Info                     models.CharacterInfo
	Attributes               models.CharacterAttributes
	Skills                   models.CharacterSkills
	validators.FormValidator `schema:"-"`
}

type userForm struct {
	Name                     string
	Password                 string
	validators.FormValidator `schema:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if userId == 0 {
		data := app.newTemplateData(r)
		app.render(w, r, http.StatusOK, "home.tmpl.html", data)
		return
	}

	characters, err := app.characters.GetAll(userId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.Characters = characters
	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = characterCreateForm{Skills: models.DefaultCharacterSkills()}

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	var form characterCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	for key, info := range form.Info.AsMap() {
		form.CheckField(validators.NotBlank(info), key, "Dieses Feld kann nicht leer sein.")
	}
	form.CheckField(validators.IsInteger(form.Info.Age), "Alter", "Dieses Feld muss eine Zahl enthalten.")
	form.CheckField(validators.InBetween(form.Info.Age, 18, 100), "Alter", "Alter muss zwischen 18 und 100 liegen.")
	form.CheckField(validators.PermittedValue(form.Info.Gender, "männlich", "weiblich"), "Geschlecht", "Geschlecht muss männlich oder weiblich sein.")

	for key, attr := range form.Attributes.AsMap() {
		if key != "BW" {
			form.CheckField(validators.PermittedValue(attr, 40, 50, 60, 70, 80), key, "Ungültiger Wert.")
		}
	}

	defaultSkills := models.DefaultCharacterSkills().AsMap()
	skillMap := form.Skills.AsMap()

	count := 0
	for skill, value := range skillMap {
		if value != defaultSkills[skill] {
			count++
			form.CheckField(validators.PermittedValue(value, 40, 50, 60, 70), skill, "Ungültiger Wert.")
		}
		if count > 9 {
			form.AddGenericError("Ungültige Fertigkeitsverteilung (mehr als 9 Fertigkeiten modifiziert).")
			break
		}
		if skill == "Finanzkraft" && value == 0 {
			form.AddFieldError("Finanzkraft", "Finanzkraft muss ein Wert zugewiesen werden.")
		}
	}

	if !validators.ValidDistribution(form.Attributes.AsMap(), []int{40, 50, 50, 50, 60, 60, 70, 80}) {
		form.AddGenericError("Ungültige Attributsverteilung.")
	}

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	_, err = app.characters.Insert(models.Character{Info: form.Info, Attributes: form.Attributes, Skills: form.Skills},
		app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {
	var form userForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "name", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.NotBlank(form.Password), "password", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.MinChars(form.Password, 8), "password", "Passwort muss mindestens 8 Zeichen lang sein.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Erfolgreich registriert! Bitte einloggen.")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userForm{}
	app.render(w, r, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	var form userForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "name", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.NotBlank(form.Password), "password", "Dieses Feld kann nicht leer sein.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Name, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddGenericError("Name und/oder Password sind falsch.")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
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
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	role, err := app.users.GetRole(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "role", role)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "role", RoleAnon)
	app.sessionManager.Put(r.Context(), "flash", "Erfolgreich ausgeloggt!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) viewCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	character, err := app.characters.Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Character = character
	app.render(w, r, http.StatusOK, "character.tmpl.html", data)
}
