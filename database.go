package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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

func (app *Application) initDB() {
	var err error
	app.Database, err = sql.Open("sqlite3", "my_database.db")
	if err != nil {
		app.ErrorLog.Fatal(err)
	}

	app.createTable()
}

func (app *Application) createTable() {
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
	_, err := app.Database.Exec(sqlStmt)
	if err != nil {
		app.ErrorLog.Fatalf("%q: %s", err, sqlStmt)
	}
}

func (app *Application) insertSeries(series Series) {
	query := "INSERT INTO Series(type, title, chapters, volumes, status, lastupdate, author, releasedate, image, rating) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, err := app.Database.Prepare(query)
	if err != nil {
		app.ErrorLog.Fatal(err)
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
		app.ErrorLog.Fatal(err)
	}

	stmt.Close()
}

func (app *Application) updateSeries(series Series) {
	query := "UPDATE Series SET type = ?, title = ?, chapters = ?, volumes = ?, status = ?, lastupdate = ?, author = ?, releasedate = ?, image = ?, rating = ? WHERE id = ?"
	stmt, err := app.Database.Prepare(query)
	if err != nil {
		app.ErrorLog.Fatal(err)
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
		app.ErrorLog.Fatal(err)
	}

	stmt.Close()
}

func (app *Application) getAllSeries() ([]Series, error) {
	query := "SELECT * FROM Series"
	rows, err := app.Database.Query(query)
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

func (app *Application) getSeriesByTitle(title string) ([]Series, error) {
	query := fmt.Sprintf("SELECT * FROM Series WHERE title LIKE '%%%s%%'", title)
	rows, err := app.Database.Query(query)
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

func (app *Application) deleteSeriesByID(id int) error {
	query := "DELETE FROM Series WHERE id = ?"
	stmt, err := app.Database.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}
