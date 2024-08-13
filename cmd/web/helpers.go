package main

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	"github.com/winik100/NoPenNoPaper/internal/validators"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	buf.WriteTo(w)
}

func (app *application) renderHtmx(w http.ResponseWriter, r *http.Request, templateName string, templateString string, data templateData) {
	t, err := template.New(templateName).Parse(templateString)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = t.ExecuteTemplate(w, templateName, data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()

	app.log.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (form *characterCreateForm) InfoChecks() {
	for key, info := range form.Info.AsMap() {
		form.CheckField(validators.NotBlank(info), key, "Dieses Feld kann nicht leer sein.")
		if key != "Geschlecht" && key != "Alter" {
			form.CheckField(validators.MaxChars(info, 50), key, "Maximal 50 Zeichen erlaubt.")
		}
	}

	form.CheckField(validators.IsInteger(form.Info.Age), "Alter", "Dieses Feld muss eine Zahl enthalten.")
	form.CheckField(validators.InBetween(form.Info.Age, 18, 100), "Alter", "Alter muss zwischen 18 und 100 liegen.")
	form.CheckField(validators.PermittedValue(form.Info.Gender, "männlich", "weiblich"), "Geschlecht", "Geschlecht muss männlich oder weiblich sein.")
}

func (form *characterCreateForm) AttributeChecks() {
	for key, attr := range form.Attributes.AsMap() {
		if key != "BW" {
			form.CheckField(validators.PermittedValue(attr, 40, 50, 60, 70, 80), key, "Ungültiger Wert.")
		}
	}

	if !validators.ValidDistribution(form.Attributes.AsMap(), []int{40, 50, 50, 50, 60, 60, 70, 80}) {
		form.AddGenericError("Ungültige Attributsverteilung.")
	}
}
