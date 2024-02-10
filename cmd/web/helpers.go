package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
	"errors"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
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

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}
	// We need to make a 'trial render' and check if that is okay before we send it to the client.
	// In order to catch runtime template errors
	buf := new(bytes.Buffer)

	// Here we write the template to the buffer instead of straight to the HTTP response writer.
	// If there is an error we can catch it and we can just then return it without sending the half HTML page.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Once no 'errors' are there we can write the HTTP status to the ResponseWriter
	w.WriteHeader(status)
	// Just writing the buffer to the HTML output.
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken: nosurf.Token(r),
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm() 
	if err != nil {
		return err 
 }

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err 
	}
	return nil 
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
