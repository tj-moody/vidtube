package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"vidtube/db"
	v1 "vidtube/routes/v1"
	"vidtube/storage"
	"vidtube/uploadingest"
)

func main() {
	ctx := context.Background()

	database, err := db.New(os.Getenv("DB_URL"))
	if err != nil {
		slog.Error("failed to open db", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	s3_endpoint := os.Getenv("S3_ENDPOINT")
	store, err := storage.New(ctx, "vidtube-videos-1", s3_endpoint)
	if err != nil {
		slog.Error("failed to init s3 client", "error", err)
		os.Exit(1)
	}

	http.HandleFunc("GET /api/v1", handleWithCors(handleAPIIndex))
	http.HandleFunc("/api/v1/videos", handleWithCors(v1.NewVideosHandler(database, store)))

	sqs_endpoint := os.Getenv("SQS_ENDPOINT")
	go func() {
		err := uploadingest.Run(ctx, database, "000000000000/vidtube-uploads", sqs_endpoint)
		if err != nil {
			slog.Error("upload event consumer failed", "error", err)
			os.Exit(1)
		}
	}()

	port := getPort()
	slog.Info("Server running at http://0.0.0.0:" + port)
	err = http.ListenAndServe("0.0.0.0:"+port, nil)
	if err != nil {
		slog.Error("server failed.", "err", err)
		os.Exit(1)
	}
}
