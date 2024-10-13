package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	webview "github.com/webview/webview_go"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Database *sql.DB
}

func main() {
	addr := flag.String("addr", ":51234", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}

	app.initDB()

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("/static/", app.staticHandler)
	mux.HandleFunc("/send", app.sendHandler)
	mux.HandleFunc("/delete", app.deleteHandler)
	mux.HandleFunc("/search", app.searchHandler)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	go func() {
		infoLog.Printf("Starting server on %s...", *addr)
		err := srv.ListenAndServe()
		errorLog.Fatal(err)
	}()

	w := webview.New(false)
	defer w.Destroy()

	w.SetTitle("Media Tracker")
	w.SetSize(1280, 720, webview.HintNone)
	w.Navigate("http://localhost:51234")

	infoLog.Println("Launching WebView Window...")
	w.Run()
}
