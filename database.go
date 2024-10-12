package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Track struct {
	Chapters int
	Volumes  int
	Status   string
}

type Series struct {
	Id     int
	Type   string
	Title  string
	Track  Track
	Author string
	Image  string
	Rating string
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "my_database.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable()
}

func createTable() {
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS Series (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
        title TEXT NOT NULL,
		chapters INTEGER NOT NULL,
		volumes INTEGER NOT NULL,
		status TEXT NOT NULL,
        author TEXT NOT NULL,
		image TEXT NOT NULL,
		rating INTEGER NOT NULL
    );`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s", err, sqlStmt)
	}
}

func insertSeries(series Series) {
	stmt, err := db.Prepare("INSERT INTO Series(type, title, chapters, volumes, status, author, image, rating) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(series.Type, series.Title, series.Track.Chapters, series.Track.Volumes, series.Track.Status, series.Author, series.Image, series.Rating)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Close()
}

func updateSeries(series Series) {
	stmt, err := db.Prepare("UPDATE Series SET type = ?, title = ?, chapters = ?, volumes = ?, status = ?, author = ?, image = ?, rating = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(series.Type, series.Title, series.Track.Chapters, series.Track.Volumes, series.Track.Status, series.Author, series.Image, series.Rating, series.Id)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Close()
}

func getAllSeries() ([]Series, error) {
	rows, err := db.Query("SELECT id, type, title, chapters, volumes, status, author, image, rating FROM Series")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seriesList []Series
	for rows.Next() {
		var series Series
		if err := rows.Scan(&series.Id, &series.Type, &series.Title, &series.Track.Chapters, &series.Track.Volumes, &series.Track.Status, &series.Author, &series.Image, &series.Rating); err != nil {
			return nil, err
		}
		seriesList = append(seriesList, series)
	}

	return seriesList, nil
}
