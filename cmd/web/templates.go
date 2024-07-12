package main

import (
	"bytes"
	"fmt"
	"net/http"
)

type templateData struct {
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{}
}

func (app *application) render(w http.ResponseWriter, r *http.Request, statusCode int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.log.Error(err.Error())
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.log.Error(err.Error())
		return
	}
	w.WriteHeader(statusCode)
	buf.WriteTo(w)
}
