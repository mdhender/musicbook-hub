// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"database/sql"
	"encoding/json"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"strconv"
)

type Book struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Instrument  string `json:"instrument"`
	Condition   string `json:"condition"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type BookListResponse struct {
	Books []Book `json:"books"`
	Count int    `json:"count"`
}

const dbPath = "books.db"

var (
	db *sql.DB
)

func booksHandler(w http.ResponseWriter, r *http.Request) {
	// log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	enableCORS(w)

	switch r.Method {
	case http.MethodGet:
		auth := isAuthenticated(r)

		filtered, err := listBooks(auth)
		if err != nil {
			// log.Printf("%s %s: failed to list books: %v\n", r.Method, r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if filtered == nil {
			filtered = []Book{}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(BookListResponse{
			Books: filtered,
			Count: len(filtered),
		})

	case http.MethodPost:
		if !isAuthenticated(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var b Book
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var err error
		b.ID, err = addBook(b)
		if err != nil {
			// log.Printf("%s %s: failed to add book: %v\n", r.Method, r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(b)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	// log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	enableCORS(w)

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !isAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	err = deleteBook(id)
	if err != nil {
		// log.Printf("%s %s: failed to delete book: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, "Book not found", http.StatusNotFound)
		//// log.Printf("%s %s: failed to delete book: %v\n", r.Method, r.URL.Path, err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func bookByIDHandler(w http.ResponseWriter, r *http.Request) {
	// log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	enableCORS(w)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := getBookByID(id)
	if err != nil {
		// log.Printf("%s %s: error finding book: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	if book == nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	if !book.Public && !isAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(book)
}

func booksExportHandler(w http.ResponseWriter, r *http.Request) {
	// log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	enableCORS(w)

	auth := isAuthenticated(r)

	books, err := listBooks(auth)
	if err != nil {
		// log.Printf("%s %s: failed to list books: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, "Failed to export books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=books-export.json")
	_ = json.NewEncoder(w).Encode(books)
}

func updateBookHandler(w http.ResponseWriter, r *http.Request) {
	// log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	enableCORS(w)

	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var updated Book
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	updated.ID = id

	if err := updateBook(updated); err != nil {
		// log.Printf("%s %s: failed to update book: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updated)
}

func loadBooks() error {
	// Check if the DB file exists
	_, err := os.Stat(dbPath)
	isNew := os.IsNotExist(err)

	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if isNew {
		if err := createSchema(); err != nil {
			log.Fatalf("failed to create schema: %v", err)
		}
		log.Println("created new database and schema")
	}

	return nil
}

func createSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS books (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		title       TEXT NOT NULL,
		author      TEXT NOT NULL DEFAULT '',
		instrument  TEXT NOT NULL DEFAULT '',
		condition   TEXT NOT NULL DEFAULT '',
		description TEXT NOT NULL DEFAULT '',
		public      INTEGER NOT NULL DEFAULT 0 CHECK (public in (0, 1))
	);
	`
	_, err := db.Exec(schema)
	return err
}

func addBook(book Book) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO books (title, author, instrument, condition, description, public)
		VALUES (?, ?, ?, ?, ?, ?)`,
		book.Title, book.Author, book.Instrument, book.Condition, book.Description, book.Public,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func updateBook(book Book) error {
	_, err := db.Exec(`
		UPDATE books
		SET title = ?, author = ?, instrument = ?, condition = ?, description = ?, public = ?
		WHERE id = ?`,
		book.Title, book.Author, book.Instrument, book.Condition, book.Description, book.Public, book.ID,
	)
	return err
}

func deleteBook(id int) error {
	_, err := db.Exec(`DELETE FROM books WHERE id = ?`, id)
	return err
}

func listBooks(isAuth bool) ([]Book, error) {
	rows, err := db.Query(`SELECT id, title, author, instrument, condition, description, public FROM books`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		var public int64
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Instrument, &b.Condition, &b.Description, &public)
		if err != nil {
			return nil, err
		}
		b.Public = public == 1
		// filter books based on public field or if authenticated
		if b.Public || isAuth {
			books = append(books, b)
		}
	}
	if books == nil {
		books = []Book{}
	}
	return books, rows.Err()
}

func getBookByID(id int) (*Book, error) {
	row := db.QueryRow(`SELECT id, title, author, instrument, condition, description, public FROM books WHERE id = ?`, id)

	var b Book
	var public int64
	err := row.Scan(&b.ID, &b.Title, &b.Author, &b.Instrument, &b.Condition, &b.Description, &public)
	if err == sql.ErrNoRows {
		return nil, nil // not found
	}
	if err != nil {
		return nil, err
	}
	b.Public = public == 1
	return &b, nil
}
