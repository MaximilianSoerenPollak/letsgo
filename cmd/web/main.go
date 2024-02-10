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
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"snippetbox.msp.net/internal/models"
)

type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel // This is the model we have imported from the 'internal' module.
	templateCache  map[string]*template.Template
	users          *models.UserModel
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// ===== Command line Arguments =====
	addr := flag.String("addr", ":4000", "HTTP Network address")
	dsn := flag.String("dsn", "web:verysecurepassword@/snippetbox?parseTime=True", "MySQL data source name")
	flag.Parse()

	// ===== Create Custom Logger =====
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// ===== Create Database Pool Connection ====
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	// ===== Initialize Application struct =====
	// initialize a new instance of the application struct with all dependencies (e.g. our logger for now)
	// De-Referencing the pointer of the Snippetmodel from the Models to have the connection pool in our App dependencies.
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:			&models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// ===== Start & Config Server and routes =====

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion: tls.VersionTLS13,
	}
	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server", slog.String("addr", ":4000"))
	// Using the logger to return / log any errors that http.ListenAndServe gives us.
	// We are also using here app.routes() in order to get the servemux etc.
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
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
