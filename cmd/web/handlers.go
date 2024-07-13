package main

import (
	"net/http"

	"github.com/winik100/NoPenNoPaper/internal/models"
)

type characterCreateForm struct {
	Name       string `form:"name"`
	Profession string `form:"profession"`
	Age        int    `form:"age"`
	Gender     string `form:"gender"`
	Residence  string `form:"residence"`
	Birthplace string `form:"birthplace"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	characters, err := app.characters.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData()
	data.Characters = characters

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

	_, err = app.characters.Insert(personalInfo, map[string]int{})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
