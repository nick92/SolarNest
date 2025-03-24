package main

import (
	"embed"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nick92/solarnest/sensors"
)

var content embed.FS

func main() {
	r := chi.NewRouter()

	// Serve static files (React build)
	r.Handle("/*", http.FileServer(http.FS(content)))

	// API route for solar/battery status
	r.Get("/api/status", func(w http.ResponseWriter, r *http.Request) {
		status := sensors.GetStatus()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})

	http.ListenAndServe(":8080", r)
}
