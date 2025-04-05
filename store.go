// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"database/sql"
	"log"
	"os"
)

const dbPath = "books.db"

var (
	db *sql.DB
)

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
