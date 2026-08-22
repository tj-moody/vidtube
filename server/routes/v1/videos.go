package v1

import (
	"context"
	"encoding/json"
	"fmt"
	uuid "github.com/google/uuid"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"vidtube/db"
)

type VideoUploadRequest struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	AuthorID int64  `json:"authorID"`
}

type VideoStore interface {
	GetVideos(pageSize, pageNumber int) ([]db.Video, error)
	UploadVideo(publicID uuid.UUID, title string, authorID int64, s3Key string) (int64, error)
}

type Presigner interface {
	GeneratePresignedUploadURL(ctx context.Context, key string, expires time.Duration) (string, error)
}

// /api/v1/videos
func NewVideosHandler(store VideoStore, presigner Presigner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			videosGetHandler(store, w, r)
		case http.MethodPost:
			videosPostHandler(store, presigner, w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func videosPostHandler(store VideoStore, presigner Presigner, w http.ResponseWriter, r *http.Request) {
	var req VideoUploadRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if req.AuthorID == 0 {
		http.Error(w, "Author is required", http.StatusBadRequest)
		return
	}

	videoID, err := uuid.NewV7() // or uuid.New() for random v4
	if err != nil {
		http.Error(w, "Failed to generate video id", http.StatusInternalServerError)
		return
	}
	s3Key := fmt.Sprintf("videos/%s/original.mp4", videoID)
	uploadURL, err := presigner.GeneratePresignedUploadURL(r.Context(), s3Key, 15*time.Minute)
	if err != nil {
		slog.Error("failed to generate presigned url", "error", err)
		http.Error(w, "Failed to generate upload url", http.StatusInternalServerError)
		return
	}

	id, err := store.UploadVideo(videoID, req.Title, req.AuthorID, s3Key)
	if err != nil {
		slog.Error("failed to add video", "error", err)
		http.Error(w, "failed to add video: ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Video added successfully",
		"id":        fmt.Sprint(id),
		"uploadURL": uploadURL,
	})
}

func videosGetHandler(store VideoStore, w http.ResponseWriter, r *http.Request) {
	countStr := r.URL.Query().Get("count")
	pageStr := r.URL.Query().Get("page")

	count, err := strconv.Atoi(countStr)
	if err != nil {
		http.Error(w, "Invalid `count`: "+countStr, http.StatusBadRequest)
	}
	if count <= 0 {
		count = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, "Invalid `page`: "+pageStr, http.StatusBadRequest)
	}
	if page < 1 {
		page = 1
	}

	videos, err := store.GetVideos(count, page)
	if err != nil {
		http.Error(w, "Internal database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(videos)
}
