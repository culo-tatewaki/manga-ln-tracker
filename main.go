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

	log.Println("Starting server on :51234...")
	err := http.ListenAndServe(":51234", mux)
	log.Fatal(err)
}
