package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
	"github.com/julienschmidt/httprouter"
	"snippetbox.msp.net/internal/models"
)

// ==== ROUTES ====
// Change signature of the home handler so it is defined against *application
// We basically make this a method of Application instead of being it's own function.
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Calling the data here in order to populate the default data (currently just the year)
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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
	// Getting new data object to populate the default values
	data := app.newTemplateData(r)
	data.Snippet = snippet
	// Execute the templae files.
	// Important to note here, anything that you pass as the final parameter to ExecuteTemplate is represented as the '.'
	// Passing in our tempalte struct here so we can access all data that is internal.
	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return 
	}
	// Get the content and stuff from the Form that we processed
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return 
	}

	// validation of the Form etc.

	fieldErrors := make(map[string]string)
	
	// check if title is < 100 char. long and not blank
	// We are using the Runecount and not len bto count the UNICODE POINTS, not the bytes.
	// e.g. -> Zoë has 3 Unicode points but 4 bytes. So we only count 3 not 4 here.
	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}
	
	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "This field cannot be blank"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors["expires"] = "This field must be equal to 1, 7 or 365"
	}
	
	// If there are any errors, just dump them in PLAIN TEXT http response (FPrint) and return from the handler
	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return 
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return 
	} 
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

