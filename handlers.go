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

func (app *Application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.InfoLog.Printf("Tried to access to an invalid path: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	seriesList, err := app.getAllSeries()
	if err != nil {
		app.InfoLog.Printf("/ failed to get the series")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFS(templateFiles, "templates/index.html"))
	tmpl.Execute(w, seriesList)
}

func (app *Application) staticHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/static/"):]
	data, err := staticFiles.ReadFile("static/" + filePath)
	if err != nil {
		app.InfoLog.Printf("static file not found: %s", filePath)
		http.NotFound(w, r)
		return
	}

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

func (app *Application) sendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.InfoLog.Printf("Tried to access /send with method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.InfoLog.Printf("Error parsing form data")
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

	app.InfoLog.Println("New Series: ", series)
	if series.Id == -1 {
		app.insertSeries(series)
	} else {
		app.updateSeries(series)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *Application) searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.InfoLog.Printf("Tried to access /search with method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.InfoLog.Println("Error parsing form data")
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	releaseDate, _ := strconv.Atoi(r.FormValue("release-date"))
	series := Series{
		Type:  r.FormValue("type"),
		Title: r.FormValue("title"),
		Track: Track{
			Status: r.FormValue("status"),
		},
		ReleaseDate: releaseDate,
		Rating:      r.FormValue("rating"),
	}

	app.InfoLog.Println("Searching Series like: ", series)
	seriesList, err := app.getSeriesBySearch(series)
	if err != nil {
		app.InfoLog.Println("Error filtering Series by title")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFS(templateFiles, "templates/index.html"))
	tmpl.Execute(w, seriesList)
}

func (app *Application) deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", http.MethodDelete)
		app.InfoLog.Printf("Tried to access /delete with method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		app.InfoLog.Printf("Invalid ID to delete: %d", id)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = app.deleteSeriesByID(id)
	if err != nil {
		app.InfoLog.Printf("Error deleting Series with id: %d", id)
		http.Error(w, "Error deleting item", http.StatusInternalServerError)
		return
	}

	app.InfoLog.Printf("Deleted Series with the ID: %d", id)
}
