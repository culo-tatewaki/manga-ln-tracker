package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
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

	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))

	mux.Handle("/", http.StripPrefix("/", fileServer))
	mux.HandleFunc("/add", app.addHandler)
	mux.HandleFunc("/update", app.updateHandler)
	mux.HandleFunc("/delete", app.deleteHandler)
	mux.HandleFunc("/search", app.searchHandler)
	mux.HandleFunc("/getall", app.getAllHandler)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s...", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
