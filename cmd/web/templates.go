package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/winik100/NoPenNoPaper/internal/models"
	"github.com/winik100/NoPenNoPaper/ui"
)

type templateData struct {
	Characters      []models.Character
	Character       models.Character
	User            models.User
	Form            any
	CSRFToken       string
	Flash           string
	IsAuthenticated bool
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CSRFToken:       nosurf.Token(r),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

func half(value int) int {
	return value / 2
}

func fifth(value int) int {
	return value / 5
}

var funcs = template.FuncMap{
	"half":  half,
	"fifth": fifth,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}
		ts, err := template.New(name).Funcs(funcs).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
