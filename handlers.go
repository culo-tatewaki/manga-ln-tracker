package main

import (
	"embed"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	seriesList, err := getAllSeries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFS(templateFiles, "templates/index.html"))
	tmpl.Execute(w, seriesList)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/static/"):]
	data, err := staticFiles.ReadFile("static/" + filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Set content type based on file extension
	switch {
	case filePath[len(filePath)-4:] == ".css":
		w.Header().Set("Content-Type", "text/css")
	case filePath[len(filePath)-3:] == ".js":
		w.Header().Set("Content-Type", "application/javascript")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.Write(data)
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
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
	releaseDate, _ := strconv.Atoi(r.FormValue("release-date"))
	series := Series{
		Id:    id,
		Type:  r.FormValue("type"),
		Title: r.FormValue("title"),
		Track: Track{
			Chapters:   chapters,
			Volumes:    volumes,
			Status:     r.FormValue("status"),
			LastUpdate: time.Now(),
		},
		Author:      r.FormValue("author"),
		ReleaseDate: releaseDate,
		Image:       r.FormValue("image"),
		Rating:      r.FormValue("rating"),
	}

	if series.Id == -1 {
		insertSeries(series)
	} else {
		updateSeries(series)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	search := r.FormValue("search")
	seriesList, err := getSeriesByTitle(search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFS(templateFiles, "templates/index.html"))
	tmpl.Execute(w, seriesList)
}
