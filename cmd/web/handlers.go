package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
	
	"snippetbox.msp.net/internal/validator"
	"snippetbox.msp.net/internal/models"

	"github.com/julienschmidt/httprouter"
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
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

type snippetCreateForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	validator.Validator `form:"-"`
}


func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return 
	}



	// validation of the Form etc.
	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return 
	}
	// check if title is < 100 char. long and not blank
	// We are using the Runecount and not len bto count the UNICODE POINTS, not the bytes.
	// e.g. -> Zoë has 3 Unicode points but 4 bytes. So we only count 3 not 4 here.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be longer than 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must be 1, 7 or 365")

	
	// If there are any validation errors, then re-display the create.tmpl template,
    // passing in the snippetCreateForm instance as dynamic data in the Form 
    // field. Note that we use the HTTP status code 422 Unprocessable Entity 
    // when sending the response to indicate that there was a validation error.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form 
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return 
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return 
	} 
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}


