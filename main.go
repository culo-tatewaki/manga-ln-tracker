package main

import (
	"log"
	"net/http"
)

func main() {
	initDB()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/send", sendHandler)

	println("Starting server on port :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
