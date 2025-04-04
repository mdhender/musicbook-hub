// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte("super-secret-change-this") // replace with env var in production
var hexSecret = hex.EncodeToString(jwtSecret)      // for logging, which is a security risk
var validUUIDs []uuid.UUID

func loadMagicKeys(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var keys []string
	if err := json.Unmarshal(data, &keys); err != nil {
		return err
	}

	for _, k := range keys {
		u, err := uuid.Parse(k)
		if err == nil {
			validUUIDs = append(validUUIDs, u)
		}
	}

	// Create a new secret by hashing jwtSecret and all the keys
	hasher := sha256.New()
	hasher.Write(jwtSecret)
	for _, k := range keys {
		hasher.Write([]byte(k))
	}
	jwtSecret = hasher.Sum(nil)

	// Convert to hex string for readability in logs
	hexSecret = hex.EncodeToString(jwtSecret)

	return nil
}

func isValidMagicUUID(id uuid.UUID) bool {
	for _, k := range validUUIDs {
		if k == id {
			return true
		}
	}
	return false
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	enableCORS(w)

	idStr := r.PathValue("uuid")
	uid, err := uuid.Parse(idStr)
	if err != nil || !isValidMagicUUID(uid) {
		http.Error(w, "Invalid or unauthorized UUID", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": uid.String(),
		"exp": time.Now().Add(14 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Token generation error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	enableCORS(w)

	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Could not extract claims", http.StatusUnauthorized)
		return
	}

	// Example: returning only UUID from 'sub'
	json.NewEncoder(w).Encode(map[string]any{
		"user": claims["sub"],
		"exp":  claims["exp"],
	})
}
