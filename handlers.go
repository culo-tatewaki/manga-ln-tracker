package main

import (
	"html/template"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	books, err := getAllBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse and execute the template
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, books)
}
