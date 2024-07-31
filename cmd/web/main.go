package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/schema"
	"github.com/winik100/NoPenNoPaper/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	log            *slog.Logger
	characters     models.CharacterModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
	formDecoder    *schema.Decoder
}

func main() {

	port := flag.String("port", ":8080", "HTTP Port")
	dsn := flag.String("dsn", "web:mellon@/NoPenNoPaper", "MySQL Data Source Name")
	//dsn := flag.String("dsn", "web:testpwweb@tcp(localhost:3307)/NoPenNoPaper", "MySQL Data Source Name")
	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	cache, err := newTemplateCache()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	formDecoder := schema.NewDecoder()
	formDecoder.IgnoreUnknownKeys(true)

	app := &application{
		log:            log,
		characters:     &models.CharacterModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  cache,
		sessionManager: sessionManager,
		formDecoder:    formDecoder,
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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
