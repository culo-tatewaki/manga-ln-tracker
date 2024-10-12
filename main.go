package main

import (
	"log"
	"net/http"
)

func main() {
	initDB()
	insertBook(Book{"toaru", "kamachi", 1, 3, []byte{1, 2}})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)

	println("Starting server on port :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
