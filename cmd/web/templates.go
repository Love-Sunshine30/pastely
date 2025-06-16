package main

import (
	"html/template"
	"path/filepath"
	"time"

	"al.imran.pastely/internal/models"
)

// create humanDate() function that returns a humanize dateTime string
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 03:04 PM")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// Create a templateData to hold all the dynamic data that we want to render on the page
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initializing a new new map
	cache := map[string]*template.Template{}

	// get all the filepath of 'page' template
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	// loop through the pages one-by-one
	for _, page := range pages {
		name := filepath.Base(page)

		// Create a new template set and parse the base template
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Add the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	// return the map
	return cache, nil
}
