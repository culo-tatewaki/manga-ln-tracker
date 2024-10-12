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

	log.Println("Starting server on :50001...")
	err := http.ListenAndServe(":50001", mux)
	log.Fatal(err)
}
