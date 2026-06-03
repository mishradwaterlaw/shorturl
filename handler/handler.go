package handler

import (
	"crypto/rand"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mishradwaterlaw/shorturl/store"
)

type Handler struct {
	store store.Store
	logCh chan string // This field stores the channel so handler methods can access it to send data.
}

// NewHandler now accepts logCh as a dependency, injecting the communications channel from main.go.
func NewHandler(store store.Store, logCh chan string) *Handler {
	return &Handler{
		store: store, // Assignments map the incoming interface and channel into our internal struct fields.
		logCh: logCh, // The struct now safely holds a reference to the shared channel.
	}
}

// 1. Define the incoming request shape
type ShortenRequest struct {
	URL string `json:"url"`
}

// 2. Define the outgoing response shape
type ShortenResponse struct {
	Code string `json:"code"`
}

// HandleShorten handles POST requests to create a short URL
func (h *Handler) HandleShorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest

	// Decode the JSON body into our struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	code, err := generateCode()
	if err != nil {
		http.Error(w, "Failed to generate short code", http.StatusInternalServerError)
		return
	}

	// Save it to our store layer
	h.store.Save(code, req.URL)

	// Set header and return the response back as JSON
	w.Header().Set("Content-Type", "application/json")
	response := ShortenResponse{Code: code}
	json.NewEncoder(w).Encode(response)
}

// HandleRedirect handles GET requests to forward users to the original URL
func (h *Handler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	// r.URL.Path gives us something like "/abcde123", so we trim the leading "/"
	code := strings.TrimPrefix(r.URL.Path, "/")

	if code == "" {
		http.Error(w, "Missing short code", http.StatusBadRequest)
		return
	}

	// Lookup the original URL using our store
	original, err := h.store.Get(code)
	if err != nil {
		// If there's an error (e.g., code not found), return a standard 404
		http.NotFound(w, r)
		return
	}

	// The `<-` operator sends data (the short code) into the channel.
	// Because logCh is a buffered channel, this send is non-blocking (instantaneous)
	// unless the buffer fills past 100 pending messages.
	h.logCh <- code

	// If found, redirect the browser to the original URL
	http.Redirect(w, r, original, http.StatusMovedPermanently)
}

func generateCode() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	bytes := make([]byte, 8)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}

	return string(bytes), nil
}
