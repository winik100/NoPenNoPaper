package main

import (
	"crypto/tls"
	"flag"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/winik100/NoPenNoPaper/ui"
)

type application struct {
	log            *slog.Logger
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {

	port := flag.String("port", ":8080", "HTTP Port")
	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cache, err := newTemplateCache()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		log:            log,
		templateCache:  cache,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
	}

	server := http.Server{
		Addr:         *port,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(log.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.log.Info("starting server", slog.String("port", *port))
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	app.log.Error(err.Error())
	os.Exit(1)
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
		ts, err := template.New(name).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
