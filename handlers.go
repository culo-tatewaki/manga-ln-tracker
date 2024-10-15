package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func (app *Application) addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.InfoLog.Printf("Tried to access /add with method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.InfoLog.Println("Failed to read the /add request body")
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var series Series
	err = json.Unmarshal(body, &series)
	if err != nil {
		app.InfoLog.Println("Failed to parse the /add request body into a go object")
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

	id, err := app.insertSeries(series)
	series.Id = id
	app.InfoLog.Println("Adding Series: ", series)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(series)
}

func (app *Application) updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.Header().Set("Allow", http.MethodPut)
		app.InfoLog.Printf("Tried to access /update with method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.InfoLog.Println("Failed to read the /update request body")
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var series Series
	err = json.Unmarshal(body, &series)
	if err != nil {
		app.InfoLog.Println("Failed to parse the /update request body into a go object")
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

	app.InfoLog.Println("updating Series: ", series)
	app.updateSeries(series)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(series)
}

func (app *Application) getAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.InfoLog.Printf("Tried to access /getall with method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	seriesList, err := app.getAllSeries()
	if err != nil {
		app.InfoLog.Println("Error getting all series")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(seriesList)
}

func (app *Application) searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.InfoLog.Printf("Tried to access /search with method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.InfoLog.Println("Failed to read the /search request body")
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var series Series
	err = json.Unmarshal(body, &series)
	if err != nil {
		app.InfoLog.Println("Failed to parse the /search request body into a go object")
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

	app.InfoLog.Println("Searching Series like: ", series)
	seriesList, err := app.getSeriesBySearch(series)
	if err != nil {
		app.InfoLog.Println("Error filtering Series by title")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(seriesList)
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
	w.WriteHeader(http.StatusNoContent)
}
