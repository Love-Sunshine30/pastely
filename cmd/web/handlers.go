package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"al.imran.pastely/internal/models"
	"al.imran.pastely/internal/validator"
	"github.com/julienschmidt/httprouter"
)

// creating a struch to hold the user sign up form and any error
type userSignUpForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// creating a struct to hold the snippet create and any error that user may input
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
	}
	// store all the data to a data variable
	// first add the current year w/newTemplateData
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// use the render helper
	app.render(w, http.StatusOK, "home.tmpl.html", data)

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// getting the parameter from request context
	params := httprouter.ParamsFromContext(r.Context())

	//Getting the id from params
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Retriving a specific snippet when requested
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// storing all the templateData to data variable
	data := app.newTemplateData(r)
	data.Snippet = snippet

	// render the page
	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// declare an instance of createSnippetForm struct
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// validate the form data
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxCharCount(form.Title, 100), "title", "This field cannot conatn more than 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal to 1, 7 or 365")

	// if validation fails, re-render the form with error message
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}
	// insert the snippet data to our db
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
	}

	// adding the flash message to the session data
	app.sessionManager.Put(r.Context(), "flash", "Snippet Created Sucessfully!")

	// Redirect to the created snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// Handler for user sign up form
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignUpForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

// Handler for signing up new user
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Declare a zero-valued instance of the signupForm struc
	var form userSignUpForm

	// parse the form data to the form struct
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validate the form data
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank!")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank!")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank!")
	form.CheckField(validator.MinCharCount(form.Password, 8), "password", "Password must be at least 8 character long.")

	// if the form has any error then re-display the form with 422 status code
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}
	// Insert the new user to the database
	err = app.user.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		// if error is ErrorDuplicateEmail then we re-display the form with a message
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFiledError("email", "This Email is already used!")

			// Re-display the form
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "You've signed up successfully!")
	// Redirect to the log in page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// Handler for login form
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {

}

// Handler for logining in a user
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {

}

// Handler for logingin out a user
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {

}
