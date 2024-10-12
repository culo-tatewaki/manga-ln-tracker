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

func sendHandler(w http.ResponseWriter, r *http.Request) {
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
	id, _ := strconv.Atoi(r.FormValue("id"))
	book := Book{
		Id:     id,
		Type:   r.FormValue("type"),
		Series: r.FormValue("series"),
		Volume: volume,
		Author: r.FormValue("author"),
		Image:  r.FormValue("image"),
		Rating: r.FormValue("rating"),
	}

	if book.Id == -1 {
		insertBook(book)
		fmt.Println("insert")

	} else {
		fmt.Println("modify")
		modifyBook(book)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
