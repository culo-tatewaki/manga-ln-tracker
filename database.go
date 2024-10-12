package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Book struct {
	Series string
	Volume int
	Author string
	Image  string
	Rating int
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "my_database.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create the table if it doesn't exist
	createTable()
}

func createTable() {
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS books (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        series TEXT NOT NULL,
		volume INTEGER NOT NULL,
        author TEXT NOT NULL,
		image BLOB NOT NULL,
		rating INTEGER
    );`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s", err, sqlStmt)
	}
}

func insertBook(book Book) {
	stmt, err := db.Prepare("INSERT INTO books(series, volume, author, image, rating) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(book.Series, book.Volume, book.Author, book.Image, book.Rating)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Close()
}

func getAllBooks() ([]Book, error) {
	rows, err := db.Query("SELECT series, volume, author, image, rating FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.Series, &book.Volume, &book.Author, &book.Image, &book.Rating); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}
