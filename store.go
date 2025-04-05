// store.go
// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"database/sql"
	"embed"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
)

const dbPath = "books.db"

var (
	db *sql.DB
)

func openStore() error {
	// Check if database file exists before opening
	if _, err := os.Stat(dbPath); err == nil {
		// Database exists, create a backup
		log.Printf("backing up database file\n")

		fileInfo, err := os.Stat(dbPath)
		if err != nil {
			log.Fatalf("error: failed to get file info for backup: %v", err)
		}
		// Format timestamp in UTC as YYYYMMDDHH24MISS
		modTime := fileInfo.ModTime().UTC()
		timestamp := modTime.Format("20060102150405")
		backupPath := dbPath + "," + timestamp

		// Copy the database file to create backup
		srcFile, err := os.Open(dbPath)
		if err != nil {
			log.Fatalf("error: failed to open source file for backup: %v", err)
		}
		defer srcFile.Close()

		// todo: don't overwrite existing backup?
		dstFile, err := os.Create(backupPath)
		if err != nil {
			log.Fatalf("error: failed to create backup file: %v", err)
		}
		defer dstFile.Close()

		if _, err = io.Copy(dstFile, srcFile); err != nil {
			log.Fatalf("error: failed to copy database for backup: %v", err)
		}
		log.Printf("created database backup: %s", backupPath)
	}

	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	// Apply any pending migrations
	if err := migrateSchema(); err != nil {
		log.Fatalf("failed to migrate schema: %v", err)
	}

	return nil
}

var (
	//go:embed migrations/*.sql
	migrationsFS embed.FS
)

func migrateSchema() error {
	// Create migrations table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			filename TEXT NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("failed to create migrations table: %v", err)
		return err
	}

	// Get list of migration files
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		log.Fatalf("failed to read migrations directory: %v", err)
		return err
	}

	// Sort migration files by name (timestamp prefix ensures correct order)
	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sql" {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Begin transaction for all migrations
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check which migrations have already been applied
	for _, filename := range migrationFiles {
		var exists bool
		err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM migrations WHERE filename = ?)", filename).Scan(&exists)
		if err != nil {
			log.Fatalf("failed to check if migration exists: %v", err)
			return err
		}

		// Skip if migration has already been applied
		if exists {
			log.Printf("Migration %s already applied, skipping", filename)
			continue
		}

		// Read and apply the migration file
		migrationSQL, err := migrationsFS.ReadFile(filepath.Join("migrations", filename))
		if err != nil {
			log.Fatalf("failed to read migration file %s: %v", filename, err)
			return err
		}

		// Execute the migration
		_, err = tx.Exec(string(migrationSQL))
		if err != nil {
			log.Fatalf("failed to apply migration %s: %v", filename, err)
			return err
		}

		// Record the migration as applied
		_, err = tx.Exec("INSERT INTO migrations (filename) VALUES (?)", filename)
		if err != nil {
			log.Fatalf("failed to record migration %s: %v", filename, err)
			return err
		}

		log.Printf("Successfully applied migration: %s", filename)
	}

	// Commit all migrations
	if err = tx.Commit(); err != nil {
		log.Fatalf("failed to commit migrations: %v", err)
		return err
	}

	log.Println("Schema migration completed successfully")
	return nil
}

func addBook(book Book) (int64, error) {
	result, err := db.Exec(`
	INSERT INTO books (title, author, condition, format, description, instrument, public)
	VALUES (?, ?, ?, ?, ?, ?, ?)`,
		book.Title, book.Author, book.Condition, book.Format, book.Description, book.Instrument, book.Public,
	)

	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func updateBook(book Book) error {
	// log.Printf("updateBook: %+v\n", book)
	_, err := db.Exec(`
		UPDATE books
		SET title = ?, author = ?, condition = ?, format = ?, description = ?, instrument = ?, public = ?
		WHERE id = ?`,
		book.Title, book.Author, book.Condition, book.Format, book.Description, book.Instrument, book.Public, book.ID,
	)
	return err
}

func deleteBook(id int64) error {
	_, err := db.Exec(`DELETE FROM books WHERE id = ?`, id)
	return err
}

func listBooks(isAuth bool) ([]Book, error) {
	rows, err := db.Query(`SELECT id, title, author, condition, format, description, instrument, public, created_at, updated_at FROM books`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		var public int64
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Condition, &b.Format, &b.Description, &b.Instrument, &public, &b.CreatedAt, &b.UpdatedAt)
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
	row := db.QueryRow(`SELECT id, title, author, condition, format, description, instrument, public, created_at, updated_at FROM books WHERE id = ?`, id)

	var b Book
	var public int64
	err := row.Scan(&b.ID, &b.Title, &b.Author, &b.Condition, &b.Format, &b.Description, &b.Instrument, &public, &b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // not found
	}
	if err != nil {
		return nil, err
	}
	b.Public = public == 1
	return &b, nil
}
