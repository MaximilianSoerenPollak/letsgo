package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"snippetbox.msp.net/internal/models"
)

// ==== ROUTES ====
// Change signature of the home handler so it is defined against *application
// We basically make this a method of Application instead of being it's own function.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r) // This function replies to the request with a 404 not found error.
		return
	}
	// ts, err := template.ParseFiles("./ui/html/pages/home.tmpl")
	// if err != nil {
	// 	log.Print(err.Error())
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }
	// err = ts.Execute(w, nil)
	// if err != nil {
	// 	log.Print(err.Error())
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// }
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
	  	// Because the home handler is now a method against the application
        // struct it can access its fields, including the structured logger. We'll 
        // use this to create a log entry at Error level containing the error
        // message, also including the request method and URI as attributes to 
        // assist with debugging.
		app.serverError(w,r,err)
		http.Error(w, "Interval Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w,r,err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // using our wrapper for convenience and shorter code
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err) 
		}
		// We return here early so it does NOT show the empty snippet to the user
		return
	}
	// Plain text response body for the HTML response
	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)                         // Leeting user know what methods are allowed once we deny his one
		app.clientError(w, http.StatusMethodNotAllowed) //use the error helpers.
		return
	}
	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\n But slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return 
	}
	// Redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
