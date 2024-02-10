package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.msp.net/internal/models"
	"snippetbox.msp.net/internal/validator"

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
	// Do not need to populate Flash as we already do that via the 'newTemplateData' from helper.
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
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
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
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Name), "name", "This field can not be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field can not be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Please entere a valid Email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Your password must have at least 8 characters")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrrDuplicateEmail) {
			form.AddFieldError("email", "Email adress already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	//Flash confirmation message on the screen
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm 
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w,http.StatusBadRequest)
		return 
	}
	form.CheckField(validator.NotBlank(form.Email), "email", "This field can not be blank") 
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email adreess") 
	form.CheckField(validator.NotBlank(form.Password), "password", "This field can not be blank") 
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form 
		app.render(w,r,http.StatusUnprocessableEntity, "login.tmpl", data)
		return 
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form 
			app.render(w,r,http.StatusUnprocessableEntity, "login.tmpl", data) 
		} else {
			app.serverError(w,r,err)
		}
		return 
	}
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return 
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// We are giving them a new user authenticate Id so they can use t he same one to log back in.
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return 
	}
	// Remove the authenticatedUserID from the session -> logging user out 
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
