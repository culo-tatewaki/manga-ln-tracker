package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {
	http.HandleFunc("/", homeHandler)

	println("Starting server on port :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
