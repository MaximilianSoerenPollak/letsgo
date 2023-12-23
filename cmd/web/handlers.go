package main

import (
	"fmt"
	"net/http"
	"strconv"
	"html/template"
	"log"
)

// ==== ROUTES ====
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r) // This function replies to the request with a 404 not found error.
		return
	}
	ts, err := template.ParseFiles("./ui/html/pages/home.tmpl.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}
	ts, err = template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Interval Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
			http.NotFound(w, r)
			return
	}
	
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) // Leeting user know what methods are allowed once we deny his one
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed) // Letting the user know that the used method is not correct
		return
	}
	w.Write([]byte("Create a new Snippet\n"))
}
