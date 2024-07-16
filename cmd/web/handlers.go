package main

import (
	"errors"
	"net/http"

	"github.com/winik100/NoPenNoPaper/internal/models"
	"github.com/winik100/NoPenNoPaper/internal/validators"
)

type characterCreateForm struct {
	Name       string `form:"name"`
	Profession string `form:"profession"`
	Age        string `form:"age"`
	Gender     string `form:"gender"`
	Residence  string `form:"residence"`
	Birthplace string `form:"birthplace"`

	ST int `form:"st"`
	GE int `form:"ge"`
	MA int `form:"ma"`
	KO int `form:"ko"`
	ER int `form:"er"`
	BI int `form:"bi"`
	GR int `form:"gr"`
	IN int `form:"in"`
	BW int

	validators.FormValidator `form:"-"`
}

func newCharacterCreateForm() characterCreateForm {
	return characterCreateForm{
		BW: 8,
	}
}

type userForm struct {
	Name                     string `form:"name"`
	Password                 string `form:"password"`
	validators.FormValidator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = characterCreateForm{}

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	form := newCharacterCreateForm()

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	info := models.CharacterInfo{
		Name:       form.Name,
		Profession: form.Profession,
		Age:        form.Age,
		Gender:     form.Gender,
		Residence:  form.Residence,
		Birthplace: form.Birthplace,
	}

	form.CheckField(validators.NotBlank(info.Name), "name", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.NotBlank(info.Profession), "profession", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.NotBlank(info.Age), "age", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.IsInteger(info.Age), "age", "Dieses Feld muss eine Zahl enthalten.")
	form.CheckField(validators.InBetween(info.Age, 18, 100), "age", "Alter muss zwischen 18 und 100 liegen.")
	form.CheckField(validators.NotBlank(info.Gender), "gender", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.PermittedValue(info.Gender, "männlich", "weiblich"), "gender", "Geschlecht muss männlich oder weiblich sein.")
	form.CheckField(validators.NotBlank(info.Residence), "residence", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.NotBlank(info.Birthplace), "birthplace", "Dieses Feld kann nicht leer sein.")

	attributes := models.CharacterAttributes{
		ST: form.ST,
		GE: form.GE,
		MA: form.MA,
		KO: form.KO,
		ER: form.ER,
		BI: form.BI,
		GR: form.GR,
		IN: form.IN,
		BW: form.BW,
	}

	for key, attr := range attributes.AsMap() {
		if key != "bw" {
			form.CheckField(validators.PermittedValue(attr, 40, 50, 60, 70, 80), key, "Ungültiger Wert.")
		}
	}

	if !validators.ValidAttributeDistribution(attributes) {
		form.AddGenericError("Ungültige Attributsverteilung.")
	}

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	_, err = app.characters.Insert(models.Character{
		Info:       info,
		Attributes: attributes,
	}, app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
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
	app.sessionManager.Put(r.Context(), "role", models.RoleAnon)
	app.sessionManager.Put(r.Context(), "flash", "Erfolgreich ausgeloggt!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) viewPlayer(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	characters, err := app.characters.GetAll(userId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Characters = characters
	app.render(w, r, http.StatusOK, "player.tmpl.html", data)
}
