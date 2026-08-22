package main

import (
	"os"
	"log/slog"
	"context"
	"net/http"
	"vidtube/db"
	"vidtube/storage"
	v1 "vidtube/routes/v1"
)

func main() {
	ctx := context.Background()

	err := db.InitDB()
	if err != nil {
		slog.Error("failed to open db.", "err", err)
		os.Exit(1)
	}
	defer db.Instance.Close()

	if err := storage.Init(ctx, "vidtube-videos-1"); err != nil {
		slog.Error("failed to init s3 client", "error", err)
		os.Exit(1)
	}

	http.HandleFunc("/api/v1/videos", handleWithCors(v1.VideosHandler))

	port := getPort()
	slog.Info("Server running at http://0.0.0.0:" + port)
	err = http.ListenAndServe("0.0.0.0:"+port, nil)
	if err != nil {
		slog.Error("server failed.", "err", err)
		os.Exit(1)
	}
}
