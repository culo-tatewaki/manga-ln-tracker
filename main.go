package main

import (
	"log"
	"net/http"
)

func main() {
	initDB()
	insertBook(Book{"toaru", 1, "kamachi", "http:/7asd", 3})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/add", addHandler)

	println("Starting server on port :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
