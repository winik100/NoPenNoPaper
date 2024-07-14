package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/winik100/NoPenNoPaper/internal/models"
)

type templateData struct {
	Characters []models.Character
	Form       any
	CSRFToken  string
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CSRFToken: nosurf.Token(r),
	}
}

func (app *application) render(w http.ResponseWriter, r *http.Request, statusCode int, page string, data templateData) {
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
	w.WriteHeader(statusCode)
	buf.WriteTo(w)
}
