package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Track struct {
	Chapters   int       `json:"chapters"`
	Volumes    int       `json:"volumes"`
	Status     string    `json:"status"`
	LastUpdate time.Time `json:"lastUpdate"`
}

type Series struct {
	Id          int64  `json:"id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Track       Track  `json:"track"`
	Author      string `json:"author"`
	ReleaseDate int    `json:"releaseDate"`
	Image       string `json:"image"`
	Rating      string `json:"rating"`
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

func (app *Application) insertSeries(series Series) (int64, error) {
	query := "INSERT INTO Series(type, title, chapters, volumes, status, lastupdate, author, releasedate, image, rating) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, err := app.Database.Prepare(query)
	if err != nil {
		app.ErrorLog.Fatal(err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
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
	)
	if err != nil {
		app.ErrorLog.Fatal(err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		app.ErrorLog.Fatal(err)
		return 0, err
	}

	return id, nil
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

	var seriesList []Series = []Series{}
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

func (app *Application) getSeriesBySearch(series Series) ([]Series, error) {
	query := "SELECT * FROM Series WHERE 1=1"
	var args []interface{}

	if series.Type != "" {
		query += " AND type = ?"
		args = append(args, series.Type)
	}
	if series.Title != "" {
		query += " AND title LIKE ?"
		args = append(args, "%"+series.Title+"%")
	}
	if series.Track.Status != "" {
		query += " AND status = ?"
		args = append(args, series.Track.Status)
	}
	if series.ReleaseDate != 0 {
		query += " AND releasedate = ?"
		args = append(args, series.ReleaseDate)
	}
	if series.Rating != "" {
		query += " AND rating = ?"
		args = append(args, series.Rating)
	}

	rows, err := app.Database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seriesList []Series = []Series{}
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
