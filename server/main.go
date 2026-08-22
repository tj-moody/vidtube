package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"vidtube/db"
	v1 "vidtube/routes/v1"
	"vidtube/storage"
)

const defaultS3Endpoint = "http://localhost:4566" // LocalStack

func main() {
	ctx := context.Background()

	database, err := db.New(os.Getenv("DB_URL"))
	if err != nil {
		slog.Error("failed to open db", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	endpoint := os.Getenv("S3_ENDPOINT")
	if endpoint == "" {
		endpoint = defaultS3Endpoint
	}
	store, err := storage.New(ctx, "vidtube-videos-1", endpoint)
	if err != nil {
		slog.Error("failed to init s3 client", "error", err)
		os.Exit(1)
	}

	http.HandleFunc("/api/v1/videos", handleWithCors(v1.NewVideosHandler(database, store)))

	port := getPort()
	slog.Info("Server running at http://0.0.0.0:" + port)
	err = http.ListenAndServe("0.0.0.0:"+port, nil)
	if err != nil {
		slog.Error("server failed.", "err", err)
		os.Exit(1)
	}
}
