// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Book struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	Instrument string `json:"instrument"`
	Condition  string `json:"condition"`
	Public     bool   `json:"public"` // ✅ New field
}

type BookListResponse struct {
	Books []Book `json:"books"`
	Count int    `json:"count"`
}

const booksFile = "books.json"

var (
	books          = []Book{}
	nextID         = 1
	booksMutex     = &sync.Mutex{}
	booksFileMutex = &sync.Mutex{}
)

func booksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	enableCORS(w)

	switch r.Method {
	case http.MethodGet:
		booksMutex.Lock()
		defer booksMutex.Unlock()

		var filtered []Book
		auth := isAuthenticated(r)

		for _, b := range books {
			if b.Public || auth {
				filtered = append(filtered, b)
			}
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
		var b Book
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		booksMutex.Lock()
		defer booksMutex.Unlock()
		b.ID = nextID
		nextID++
		books = append(books, b)
		nextID++
		_ = saveBooks() // ✅
		books = append(books, b)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(b)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	enableCORS(w)

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	booksMutex.Lock()
	defer booksMutex.Unlock()

	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			_ = saveBooks() // ✅
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}

func updateBookHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	enableCORS(w)

	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var update struct {
		Public *bool `json:"public"`
	}

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if update.Public == nil {
		http.Error(w, "Missing 'public' field", http.StatusBadRequest)
		return
	}

	booksMutex.Lock()
	defer booksMutex.Unlock()

	for i, b := range books {
		if b.ID == id {
			books[i].Public = *update.Public
			_ = saveBooks() // ✅

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(books[i])
			return
		}
	}

	log.Printf("%s %s: book %q not found\n", r.Method, r.URL.Path, r.PathValue("id"))

	http.Error(w, "Book not found", http.StatusNotFound)
}

func loadBooks() error {
	booksFileMutex.Lock()
	defer booksFileMutex.Unlock()

	file, err := os.ReadFile(booksFile)
	if err != nil {
		if os.IsNotExist(err) {
			books = []Book{}
			return nil
		}
		return err
	}

	tmp := []Book{}
	err = json.Unmarshal(file, &tmp)
	if err != nil {
		return err
	}
	books = tmp

	// Recalculate nextID based on existing books
	nextID = 1
	for _, b := range books {
		if b.ID >= nextID {
			nextID = b.ID + 1
		}
	}

	return nil
}

func saveBooks() error {
	booksFileMutex.Lock()
	defer booksFileMutex.Unlock()

	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return err
	}
	log.Printf("%s: 💾 saving %d books\n", booksFile, len(books))
	return os.WriteFile(booksFile, data, 0644)
}
