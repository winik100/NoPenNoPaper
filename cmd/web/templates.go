package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"github.com/justinas/nosurf"
	"github.com/winik100/NoPenNoPaper/internal/core"
	"github.com/winik100/NoPenNoPaper/ui"
)

type templateData struct {
	Characters      []core.Character
	Character       core.Character
	User            core.User
	Form            any
	AdditionalData  any
	CSRFToken       string
	Flash           string
	IsAuthenticated bool
	IsAuthorized    bool
}

func (app *application) newTemplateData(r *http.Request) templateData {
	user := core.User{
		ID:   app.sessionManager.GetInt(r.Context(), authenticatedUserIdKey),
		Name: app.sessionManager.GetString(r.Context(), authenticatedUserNameKey)}
	return templateData{
		User:            user,
		CSRFToken:       nosurf.Token(r),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		IsAuthorized:    app.isAuthorized(r),
	}
}

func half(value int) int {
	res := value / 2
	if res == 0 {
		return 1
	}
	return res
}

func fifth(value int) int {
	res := value / 5
	if res == 0 {
		return 1
	}
	return res
}

func contains(skills []string, skill string) bool {
	return slices.Contains(skills, skill)
}

func trim(s string) string {
	return strings.Join(strings.Split(s, " "), "")
}

var funcs = template.FuncMap{
	"half":     half,
	"fifth":    fifth,
	"contains": contains,
	"trim":     trim,
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
