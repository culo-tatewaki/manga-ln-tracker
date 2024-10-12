package main

import (
	"log"
	"net/http"
)

func main() {
	initDB()
	//insertBook(Book{1, "Manga", "Mayonaka Heart Tune", 6, "Igarashi Masakuni", "https://meo.comick.pictures/l6X18Y.jpg", "⭐⭐⭐⭐⭐"})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/send", sendHandler)

	println("Starting server on port :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
