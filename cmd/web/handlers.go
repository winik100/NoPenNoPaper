package main

import (
	"net/http"
	"strconv"

	"github.com/winik100/NoPenNoPaper/internal/models"
	"github.com/winik100/NoPenNoPaper/internal/validators"
)

type characterCreateForm struct {
	Name       string `form:"name"`
	Profession string `form:"profession"`
	Age        int    `form:"age"`
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
	BW int `form:"bw"`

	validators.FormValidator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Form = characterCreateForm{}

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	var form characterCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	personalInfo := models.PersonalInfo{
		Name:       form.Name,
		Profession: form.Profession,
		Age:        form.Age,
		Gender:     form.Gender,
		Residence:  form.Residence,
		Birthplace: form.Birthplace}

	form.CheckField(validators.NotBlank(personalInfo.Name), "name", "Dieses Feld kann nicht leer sein.")

	form.CheckField(validators.NotBlank(personalInfo.Profession), "profession", "Dieses Feld kann nicht leer sein.")

	form.CheckField(validators.NotBlank(strconv.Itoa(personalInfo.Age)), "age", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.InBetween(personalInfo.Age, 18, 100), "age", "Alter muss zwischen 18 und 100 liegen.")

	form.CheckField(validators.NotBlank(personalInfo.Gender), "gender", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.PermittedValue(personalInfo.Gender, "männlich", "weiblich"), "gender", "Geschlecht muss männlich oder weiblich sein.")

	form.CheckField(validators.NotBlank(personalInfo.Residence), "residence", "Dieses Feld kann nicht leer sein.")

	form.CheckField(validators.NotBlank(personalInfo.Birthplace), "birthplace", "Dieses Feld kann nicht leer sein.")

	attributes := models.Attributes{
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

	form.CheckField(validators.PermittedValue(attributes.ST, 10, 20, 30, 40, 50, 60, 70, 80, 90), "st", "Ungültiger Wert.")
	form.CheckField(validators.PermittedValue(attributes.GE, 10, 20, 30, 40, 50, 60, 70, 80, 90), "ge", "Ungültiger Wert.")
	form.CheckField(validators.PermittedValue(attributes.MA, 10, 20, 30, 40, 50, 60, 70, 80, 90), "ma", "Ungültiger Wert.")
	form.CheckField(validators.PermittedValue(attributes.KO, 10, 20, 30, 40, 50, 60, 70, 80, 90), "ko", "Ungültiger Wert.")
	form.CheckField(validators.PermittedValue(attributes.ER, 10, 20, 30, 40, 50, 60, 70, 80, 90), "er", "Ungültiger Wert.")
	form.CheckField(validators.PermittedValue(attributes.BI, 10, 20, 30, 40, 50, 60, 70, 80, 90), "bi", "Ungültiger Wert.")
	form.CheckField(validators.PermittedValue(attributes.GR, 10, 20, 30, 40, 50, 60, 70, 80, 90), "gr", "Ungültiger Wert.")
	form.CheckField(validators.PermittedValue(attributes.IN, 10, 20, 30, 40, 50, 60, 70, 80, 90), "in", "Ungültiger Wert.")
	form.CheckField(validators.PermittedValue(attributes.BW, 10, 20, 30, 40, 50, 60, 70, 80, 90), "bw", "Ungültiger Wert.")

	if !form.Valid() {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	_, err = app.characters.Insert(personalInfo, attributes)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
