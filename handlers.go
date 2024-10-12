package main

import (
	"html/template"
	"net/http"
	"strconv"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	seriesList, err := getAllSeries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, seriesList)
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

	chapters, _ := strconv.Atoi(r.FormValue("chapters"))
	volumes, _ := strconv.Atoi(r.FormValue("volumes"))
	id, _ := strconv.Atoi(r.FormValue("id"))
	series := Series{
		Id:    id,
		Type:  r.FormValue("type"),
		Title: r.FormValue("title"),
		Track: Track{
			Chapters: chapters,
			Volumes:  volumes,
		},
		Author: r.FormValue("author"),
		Image:  r.FormValue("image"),
		Rating: r.FormValue("rating"),
	}

	// fmt.Println(series)
	if series.Id == -1 {
		insertSeries(series)
	} else {
		updateSeries(series)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
