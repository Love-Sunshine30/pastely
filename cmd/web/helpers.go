package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

// checking if an incoming request is made from a authenticated user or not
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticationUserId")
}

// creating a helper that would decode the html form and put the data in repective struct fields
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// parse the form
	err := r.ParseForm()
	if err != nil {
		return err
	}
	// call the decoder
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning
		// the error.
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		// otherwise just return the error
		return err
	}
	return nil
}

// creating newTemplateData which returns a pointer to the templateData struct initialize
// with CurrentYear, any flash message and whether or not the user is authenticated
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

// Rendering the cached template pages
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// Retrive appropriate template set
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}
	// initialize a new buffer
	buf := new(bytes.Buffer)

	// write the template to a buffer, if error occurs, throw an error
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// if everything is ok, write out the template to http.ResponseWriter
	w.WriteHeader(status)
	buf.WriteTo(w)

}

// This will diaplay server error
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLogger.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// This will display a client error
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// This will display pageNotFound error
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
