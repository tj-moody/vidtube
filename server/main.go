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
	err := db.InitDB()
	if err != nil {
		slog.Error("Failed to open DB: ", err)
		os.Exit(1)
	}
	defer db.Instance.Close()


	http.HandleFunc("/api/v1/videos", handleWithCors(v1.VideosHandler))

	port := getPort()
	slog.Info("Server running at http://0.0.0.0:" + port)
	err = http.ListenAndServe("0.0.0.0:"+port, nil)
	if err != nil {
		slog.Error("Server failed: %s", err)
		os.Exit(1)
	}
}
