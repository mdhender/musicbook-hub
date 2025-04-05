package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	if pwd, err := os.Getwd(); err != nil {
		log.Fatalf("❌ failed to get working directory: %v", err)
	} else {
		log.Printf("🏠 working directory: %s\n", pwd)
	}

	if err := loadMagicKeys("magic-keys.json"); err != nil {
		log.Fatalf("failed to load magic keys: %v", err)
	}

	if err := loadBooks(); err != nil {
		log.Fatalf("%s: failed to load: %v", dbPath, err)
	} else {
		log.Printf("%s: ✅ loaded books store\n", dbPath)
	}

	mux := http.NewServeMux()

	// Magic login endpoint
	mux.HandleFunc("GET /api/login/{uuid}", loginHandler)
	mux.HandleFunc("GET /api/me", requireAuth(meHandler))

	// Books endpoints
	mux.HandleFunc("GET /api/books", getBooksHandler)                // public read
	mux.HandleFunc("POST /api/books", requireAuth(postBooksHandler)) // protected write
	mux.HandleFunc("GET /api/books/export", booksExportHandler)
	mux.HandleFunc("GET /api/books/{id}", bookByIDHandler)
	mux.HandleFunc("DELETE /api/books/{id}", requireAuth(deleteBookHandler)) // protected delete
	mux.HandleFunc("PATCH /api/books/{id}", requireAuth(updateBookHandler))

	mux.HandleFunc("OPTIONS /{rest...}", func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("%s %s: pre-flight\n", r.Method, r.URL.Path)
		enableCORS(w)
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("%s %s: not found\n", r.Method, r.URL.Path)
		enableCORS(w)
		http.Error(w, "generic not found", http.StatusNotFound)
	})

	log.Println("Listening on http://localhost:8181")
	log.Fatal(http.ListenAndServe(":8181", mux))
}
