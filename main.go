package main

import (
	"log"
	"net/http"
)

func main() {
	initDB()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/send", sendHandler)
	mux.HandleFunc("/search", searchHandler)

	println("Starting server on :8081...")
	err := http.ListenAndServe(":8081", mux)
	log.Fatal(err)
}
