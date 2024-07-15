package main

import (
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

type signupForm struct {
	ID                       int    `form:"id"`
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

	info := map[string]string{
		"name":       form.Name,
		"profession": form.Profession,
		"age":        form.Age,
		"gender":     form.Gender,
		"residence":  form.Residence,
		"birthplace": form.Birthplace}

	form.CheckField(validators.NotBlank(info["name"]), "name", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.NotBlank(info["profession"]), "profession", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.NotBlank(info["age"]), "age", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.IsInteger(info["age"]), "age", "Dieses Feld muss eine Zahl enthalten.")
	form.CheckField(validators.InBetween(info["age"], 18, 100), "age", "Alter muss zwischen 18 und 100 liegen.")
	form.CheckField(validators.NotBlank(info["gender"]), "gender", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.PermittedValue(info["gender"], "männlich", "weiblich"), "gender", "Geschlecht muss männlich oder weiblich sein.")
	form.CheckField(validators.NotBlank(info["residence"]), "residence", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.NotBlank(info["birthplace"]), "birthplace", "Dieses Feld kann nicht leer sein.")

	attributes := map[string]int{
		"st": form.ST,
		"ge": form.GE,
		"ma": form.MA,
		"ko": form.KO,
		"er": form.ER,
		"bi": form.BI,
		"gr": form.GR,
		"in": form.IN,
		"bw": form.BW,
	}

	for key, attr := range attributes {
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
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = signupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {
	var form signupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "name", "This field cannot be blank.")
	form.CheckField(validators.NotBlank(form.Password), "password", "This field cannot be blank.")
	form.CheckField(validators.MinChars(form.Password, 8), "password", "Password must contain at least 8 characters.")

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

	//app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please Login.")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
