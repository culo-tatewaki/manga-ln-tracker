package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	books, err := getAllBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, books)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	volume, _ := strconv.Atoi(r.FormValue("volume"))
	rating, _ := strconv.Atoi(r.FormValue("rating"))
	data := Book{
		Series: r.FormValue("series"),
		Volume: volume,
		Author: r.FormValue("author"),
		Image:  r.FormValue("image"),
		Rating: rating,
	}

	fmt.Println(data)
	http.Redirect(w, r, "/", http.StatusFound)
}
