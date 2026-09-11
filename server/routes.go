package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

var apiIndex = map[string]string{
	"GET /api/v1/videos":  "list ready videos (query: count, page)",
	"POST /api/v1/videos": "create video and get presigned upload url!",
}

func handleAPIIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiIndex)
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")

	w.Header().Set("Vary", "Origin")
	w.Header().Set("Vary", "Access-Control-Request-Method")
	w.Header().Set("Vary", "Access-Control-Request-Headers")
}

func handleWithCors(handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "body", r.Body)
		slog.Debug("", "head", r.Header)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler(w, r)

	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		slog.Warn("$PORT not found... defaulting to 8080")
		port = "8080"
	} else {
		slog.Info("$PORT found... assigned " + port)
	}
	return port
}
