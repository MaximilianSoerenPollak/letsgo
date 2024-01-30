package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)        // Leeting user know what methods are allowed once we deny his one
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
