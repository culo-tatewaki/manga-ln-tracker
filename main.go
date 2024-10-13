package main

import (
	"log"
	"net/http"

	webview "github.com/webview/webview_go"
)

func main() {
	initDB()

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/static/", staticHandler)
	mux.HandleFunc("/send", sendHandler)
	mux.HandleFunc("/search", searchHandler)

	go func() {
		log.Println("Starting server on :51234...")
		err := http.ListenAndServe(":51234", mux)
		log.Fatal(err)
	}()

	// Create a new webview instance
	w := webview.New(false)
	defer w.Destroy()

	w.SetTitle("Go Server with WebView")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://localhost:51234")

	w.Run()
}
