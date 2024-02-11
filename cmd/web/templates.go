package main

import (
	"html/template"
	"path/filepath"
	"time"
	"io/fs"

	"snippetbox.msp.net/internal/models"
	"snippetbox.msp.net/ui"
)

// This is defined as a template type to hold any data that we want to pass into our HTML templates.
// AS we can at most load one thing into the HTML Templates we want to circumvent this by adding all needed information
// to a template Struct.
type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
	Form        any // Holy cow and 'any' type?
	Flash       string // Used for a flash message
	IsAuthenticated bool  
	CSRFToken string
}

func newTemplateCache() (map[string]*template.Template, error) {
	//Initialize a new map to act as the cache -> Q: Is this normal to use a map as cache?
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}

func humanDate(t time.Time) string {
	if t.IsZero(){
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
