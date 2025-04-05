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
	id, err := strconv.ParseInt(idStr, 10, 64)
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
	id, err := strconv.ParseInt(idStr, 10, 64)
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
	// log.Printf("%s %s: id %q: %d\n", r.Method, r.URL.Path, idStr, id)

	var updated struct {
		Public *bool `json:"public"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	if updated.Public == nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	log.Printf("updated: %+v\n", updated)
	// fetch the book from the database
	book, err := getBookByID(id)
	if err != nil {
		log.Printf("%s %s: failed to get book: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	} else if book == nil {
		// log.Printf("%s %s: failed to find book\n", r.Method, r.URL.Path)
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	log.Printf("fetch: book %+v\n", *book)
	book.Public = *updated.Public
	log.Printf("updated: book %+v\n", *book)

	if err := updateBook(*book); err != nil {
		// log.Printf("%s %s: failed to update book: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(book)
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
	log.Printf("updateBook: %+v\n", book)
	_, err := db.Exec(`
		UPDATE books
		SET title = ?, author = ?, instrument = ?, condition = ?, description = ?, public = ?
		WHERE id = ?`,
		book.Title, book.Author, book.Instrument, book.Condition, book.Description, book.Public, book.ID,
	)
	return err
}

func deleteBook(id int64) error {
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

func getBookByID(id int64) (*Book, error) {
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

var (
	formatPickList = map[string]string{
		"sheet music":      "Single piece or short folio, typically softcover and intended for performance",
		"music book":       "Bound book of sheet music or exercises, such as method books, anthologies, or collections",
		"method book":      "Instructional book for learning an instrument, often organized by level",
		"score":            "Full musical score for ensembles, orchestras, or chamber music",
		"lead sheet":       "Single-line melody with chords, often used in jazz or pop music",
		"fake book":        "Gig-style book with hundreds of lead sheets for performance",
		"manuscript":       "Blank staff paper or note-taking formats for composers and students",
		"programming book": "Technical or instructional book focused on coding, software, or computer science topics",
		"textbook":         "Educational book intended for academic study, often includes theory and exercises",
		"reference book":   "Non-fiction book used for lookups or guidance, such as dictionaries, style guides, or API references",
	}
)

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
	CREATE TABLE IF NOT EXISTS format_picklist (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
	    format      TEXT NOT NULL,
	    description TEXT NOT NULL
	);
	`
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}
	// Populate the format_picklist table with data from formatPickList map
	for format, description := range formatPickList {
		_, err := db.Exec(`
			INSERT OR IGNORE INTO format_picklist (format, description)
			VALUES (?, ?)`,
			format, description)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
insert into format_picklist(format, description) values('sheet music',      'Single piece or short folio, typically softcover and intended for performance');
insert into format_picklist(format, description) values('music book',       'Bound book of sheet music or exercises, such as method books, anthologies, or collections');
insert into format_picklist(format, description) values('method book',      'Instructional book for learning an instrument, often organized by level');
insert into format_picklist(format, description) values('score',            'Full musical score for ensembles, orchestras, or chamber music');
insert into format_picklist(format, description) values('lead sheet',       'Single-line melody with chords, often used in jazz or pop music');
insert into format_picklist(format, description) values('fake book',        'Gig-style book with hundreds of lead sheets for performance');
insert into format_picklist(format, description) values('manuscript',       'Blank staff paper or note-taking formats for composers and students');
insert into format_picklist(format, description) values('programming book', 'Technical or instructional book focused on coding, software, or computer science topics');
insert into format_picklist(format, description) values('textbook',         'Educational book intended for academic study, often includes theory and exercises');
insert into format_picklist(format, description) values('reference book',   'Non-fiction book used for lookups or guidance, such as dictionaries, style guides, or API references');
*/
