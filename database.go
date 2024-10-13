package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Track struct {
	Chapters   int
	Volumes    int
	Status     string
	LastUpdate time.Time
}

type Series struct {
	Id          int
	Type        string
	Title       string
	Track       Track
	Author      string
	ReleaseDate int
	Image       string
	Rating      string
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
		lastupdate TEXT NOT NULL,
        author TEXT NOT NULL,
		releasedate INTEGER NOT NULL,
		image TEXT NOT NULL,
		rating INTEGER NOT NULL
    );`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s", err, sqlStmt)
	}
}

func insertSeries(series Series) {
	query := "INSERT INTO Series(type, title, chapters, volumes, status, lastupdate, author, releasedate, image, rating) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = stmt.Exec(
		series.Type,
		series.Title,
		series.Track.Chapters,
		series.Track.Volumes,
		series.Track.Status,
		series.Track.LastUpdate.Format("2006-01-02 15:04:05"),
		series.Author,
		series.ReleaseDate,
		series.Image,
		series.Rating,
	); err != nil {
		log.Fatal(err)
	}

	stmt.Close()
}

func updateSeries(series Series) {
	query := "UPDATE Series SET type = ?, title = ?, chapters = ?, volumes = ?, status = ?, lastupdate = ?, author = ?, releasedate = ?, image = ?, rating = ? WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = stmt.Exec(
		series.Type,
		series.Title,
		series.Track.Chapters,
		series.Track.Volumes,
		series.Track.Status,
		series.Track.LastUpdate.Format("2006-01-02 15:04:05"),
		series.Author,
		series.ReleaseDate,
		series.Image,
		series.Rating,
		series.Id,
	); err != nil {
		log.Fatal(err)
	}

	stmt.Close()
}

func getAllSeries() ([]Series, error) {
	query := "SELECT * FROM Series"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seriesList []Series
	for rows.Next() {
		var series Series
		var dateStr string
		if err := rows.Scan(
			&series.Id,
			&series.Type,
			&series.Title,
			&series.Track.Chapters,
			&series.Track.Volumes,
			&series.Track.Status,
			&dateStr,
			&series.Author,
			&series.ReleaseDate,
			&series.Image,
			&series.Rating,
		); err != nil {
			return nil, err
		}

		series.Track.LastUpdate, err = time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			return nil, err
		}

		seriesList = append(seriesList, series)
	}

	return seriesList, nil
}

func getSeriesByTitle(title string) ([]Series, error) {
	query := fmt.Sprintf("SELECT * FROM Series WHERE title LIKE '%%%s%%'", title)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seriesList []Series
	var dateStr string
	for rows.Next() {
		var series Series
		if err := rows.Scan(
			&series.Id,
			&series.Type,
			&series.Title,
			&series.Track.Chapters,
			&series.Track.Volumes,
			&series.Track.Status,
			&dateStr,
			&series.Author,
			&series.ReleaseDate,
			&series.Image,
			&series.Rating,
		); err != nil {
			return nil, err
		}

		series.Track.LastUpdate, err = time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			return nil, err
		}

		seriesList = append(seriesList, series)
	}

	return seriesList, nil
}
