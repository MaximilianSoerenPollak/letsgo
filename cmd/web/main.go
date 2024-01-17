package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	// ===== Command line Arguments =====
	addr := flag.String("addr", ":4000", "HTTP Network address")
	flag.Parse()

	// ===== Create Custom Logger =====
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// ===== Initialize Application struct =====
	// initialize a new instance of the application struct with all dependencies (e.g. our logger for now)
	app := &application{
		logger: logger,
	}

	// ===== Start & Config Server and routes =====


	logger.Info("starting server", slog.String("addr", ":4000"))
	// Using the logger to return / log any errors that http.ListenAndServe gives us.
	// We are also using here app.routes() in order to get the servemux etc. 
	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)

}
