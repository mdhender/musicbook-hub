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
		log.Printf("%s: failed to load: %v", booksFile, err)
		if os.IsNotExist(err) {
			log.Printf("%s: 📂 not found — creating new empty collection", booksFile)
			books = []Book{}
			if saveErr := saveBooks(); saveErr != nil {
				log.Fatalf("%s: ❌ failed to create: %v", booksFile, saveErr)
			}
		} else {
			log.Fatalf("%s: ❌ failed to load: %v", booksFile, err)
		}
	} else {
		log.Printf("%s: ✅ loaded %d books\n", booksFile, len(books))
	}

	mux := http.NewServeMux()

	// Magic login endpoint
	mux.HandleFunc("GET /api/login/{uuid}", loginHandler)
	mux.HandleFunc("GET /api/me", requireAuth(meHandler))

	// Books endpoints
	mux.HandleFunc("GET /api/books", booksHandler)                     // public read
	mux.HandleFunc("POST /api/books", requireAuth(booksHandler))       // protected write
	mux.HandleFunc("DELETE /api/books/{id}", requireAuth(bookHandler)) // protected delete
	mux.HandleFunc("PATCH /api/books/{id}", requireAuth(updateBookHandler))

	mux.HandleFunc("OPTIONS /{rest...}", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s: pre-flight\n", r.Method, r.URL.Path)
		enableCORS(w)
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s: not found\n", r.Method, r.URL.Path)
		enableCORS(w)
		http.Error(w, "generic not found", http.StatusNotFound)
	})

	log.Println("Listening on http://localhost:8181")
	log.Fatal(http.ListenAndServe(":8181", mux))
}
