package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Store shortened URLs in an in-memory map
var urlStore = make(map[string]string)

// URLRequest represents the JSON request body
type URLRequest struct {
	LongURL string `json:"long_url"`
}

// Response structure
type URLResponse struct {
	ShortURL string `json:"short_url"`
}

// Generate a random shortcode
func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// Shorten URL Handler
func shortenURL(w http.ResponseWriter, r *http.Request) {
	var req URLRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.LongURL == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Generate a shortcode
	shortCode := generateShortCode()
	urlStore[shortCode] = req.LongURL

	// Return the shortened URL
	baseURL := "https://url-shortener-djd8.onrender.com/"
	response := URLResponse{ShortURL: baseURL + shortCode}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Redirect to the original URL
func redirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortcode"]

	longURL, exists := urlStore[shortCode]
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/shorten", shortenURL).Methods("POST")
	r.HandleFunc("/{shortcode}", redirectURL).Methods("GET")

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	// Start the server on the determined port
	http.ListenAndServe(":" + port, r)
}
