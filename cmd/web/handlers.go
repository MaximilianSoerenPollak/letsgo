package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
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
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		http.Error(w, "Interval Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)                         // Leeting user know what methods are allowed once we deny his one
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed) // Letting the user know that the used method is not correct
		return
	}
	w.Write([]byte("Create a new Snippet\n"))
}
