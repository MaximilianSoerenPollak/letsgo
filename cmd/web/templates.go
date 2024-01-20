package main 

import "snippetbox.msp.net/internal/models"


//This is defined as a template type to hold any data that we want to pass into our HTML templates.
// AS we can at most load one thing into the HTML Templates we want to circumvent this by adding all needed information
// to a template Struct.
type templateData struct {
	Snippet models.Snippet
	Snippets []models.Snippet
}
