package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Book struct {
	Id     int
	Type   string
	Series string
	Volume int
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

	// Create the table if it doesn't exist
	createTable()
}

func createTable() {
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS books (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
        series TEXT NOT NULL,
		volume INTEGER NOT NULL,
        author TEXT NOT NULL,
		image TEXT NOT NULL,
		rating INTEGER NOT NULL
    );`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s", err, sqlStmt)
	}
}

func insertBook(book Book) {
	stmt, err := db.Prepare("INSERT INTO books(type, series, volume, author, image, rating) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(book.Type, book.Series, book.Volume, book.Author, book.Image, book.Rating)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Close()
}

func modifyBook(book Book) {
	stmt, err := db.Prepare("UPDATE books SET type = ?, series = ?, volume = ?, author = ?, image = ?, rating = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(book.Type, book.Series, book.Volume, book.Author, book.Image, book.Rating, book.Id)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Close()
}

func getAllBooks() ([]Book, error) {
	rows, err := db.Query("SELECT id, type, series, volume, author, image, rating FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.Id, &book.Type, &book.Series, &book.Volume, &book.Author, &book.Image, &book.Rating); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}
