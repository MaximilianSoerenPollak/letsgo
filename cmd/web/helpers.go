package main

import (
	"net/http"
)

// This will help us write a log entry at error level, including the requst that cause import
// Then it will send a generic 500 error to the client
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// THis sends a status specific response to the user with the corresponding description.
// For example we will use this later to send stuff like 400 Bad request back if the request was bad.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}


// This is a wrapper around the ´clientError'  that sends a 404.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
