package main

import (
	"html/template"
	"path/filepath"
	"snippetbox.msp.net/internal/models"
)

// This is defined as a template type to hold any data that we want to pass into our HTML templates.
// AS we can at most load one thing into the HTML Templates we want to circumvent this by adding all needed information
// to a template Struct.
type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	//Initialize a new map to act as the cache -> Q: Is this normal to use a map as cache?
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
