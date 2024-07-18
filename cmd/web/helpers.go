package main

import (
	"errors"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/gorilla/schema"
	"github.com/winik100/NoPenNoPaper/internal/models"
	"github.com/winik100/NoPenNoPaper/ui"
)

var translationKey = map[string]string{
	"Anthropologie":           "Anthropology",
	"Archäologie":             "Archaeology",
	"Autofahren":              "Driving",
	"Bibliotheksnutzung":      "LibraryResearch",
	"Buchführung":             "Accounting",
	"Charme":                  "Charme",
	"Cthulhu-Mythos":          "CthulhuMythos",
	"Einschüchtern":           "Intimidate",
	"Elektrische Reparaturen": "ElectricRepairs",
	"Erste Hilfe":             "FirstAid",
	"Finanzkraft":             "Financials",
	"Geschichte":              "History",
	"Horchen":                 "Listening",
	"Kaschieren":              "Concealing",
	"Klettern":                "Climbing",
	"Mechanische Reparaturen": "MechanicalRepairs",
	"Medizin":                 "Medicine",
	"Naturkunde":              "NaturalHistory",
	"Okkultismus":             "Occultism",
	"Orientierung":            "Orientation",
	"Psychoanalyse":           "PsychoAnalysis",
	"Psychologie":             "Psychology",
	"Rechtswesen":             "Law",
	"Reiten":                  "Horseriding",
	"Schließtechnik":          "Locks",
	"Schweres Gerät":          "HeavyMachinery",
	"Schwimmen":               "Swimming",
	"Springen":                "Jumping",
	"Spurensuche":             "Tracking",
	"Überreden":               "Persuasion",
	"Überzeugen":              "Convincing",
	"Verborgen bleiben":       "Stealth",
	"Verborgenes erkennen":    "DetectingSecrets",
	"Verkleiden":              "Disguising",
	"Werfen":                  "Throwing",
	"Werte schätzen":          "Valuation",
}

func translateToFieldName(skill string) string {
	return translationKey[skill]
}

func half(value int) int {
	return value / 2
}

func fifth(value int) int {
	return value / 5
}

var funcs = template.FuncMap{
	"half":                   half,
	"fifth":                  fifth,
	"translate":              translateToFieldName,
	"defaultCharacterSkills": models.DefaultCharacterSkills,
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

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var conversionError *schema.ConversionError

		if errors.As(err, &conversionError) {
			panic(err)
		}

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

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
