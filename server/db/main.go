package db

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Database struct {
	conn *sql.DB
}

type Video struct {
	ID       int    `json:"-"`
	PublicID string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Views    int    `json:"views"`
}

type User struct {
	ID      int64     `json:"id"`
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	Hash    string    `json:"-"`
	Created time.Time `json:"created"`
}

func New(dbURL string) (*Database, error) {
	if dbURL == "" {
		return nil, errors.New("DB_URL not set in environment")
	}
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{conn}, nil
}

// REQUIRES:
// `page_number` >= 1
// `page_size` >= 1
func (d *Database) GetVideos(page_size int, page_number int) ([]Video, error) {
	if page_number < 1 || page_size < 1 {
		return nil, errors.New("invalid page number or size")
	}

	offset := (page_number - 1) * page_size
	rows, err := d.conn.Query(
		`SELECT videos.id, videos.public_id, videos.title, users.name, videos.views
		FROM videos JOIN users ON videos.author_id = users.id
		WHERE videos.status = 'ready'
		ORDER BY videos.id DESC
		LIMIT $1 OFFSET $2`,
		page_size, offset,
	)
	if err != nil {
		slog.Error("failed to fetch videos", "error", err)
		return nil, err
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var v Video
		if err := rows.Scan(&v.ID, &v.PublicID, &v.Title, &v.Author, &v.Views); err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}
	return videos, nil
}

func (d *Database) UploadVideo(publicID uuid.UUID, title string, authorID int64, s3Key string) (int64, error) {
	if title == "" {
		return 0, errors.New("title cannot be empty")
	}
	if authorID == 0 {
		return 0, errors.New("author id cannot be empty")
	}

	var id int64
	err := d.conn.QueryRow(
		"INSERT INTO videos (public_id, title, author_id, s3_key) VALUES ($1, $2, $3, $4) RETURNING id",
		publicID,
		title,
		authorID,
		s3Key,
	).Scan(&id)
	if err != nil {
		slog.Error("failed to insert video", "error", err, "title", title, "author_id", authorID)
		return 0, err
	}

	slog.Info("successfully added video", "title", title, "author_id", authorID, "id", id)
	return id, nil
}

func (d *Database) MarkVideoReady(publicID uuid.UUID) error {
	res, err := d.conn.Exec(
		"UPDATE videos SET status = 'ready' WHERE public_id = $1 AND status IN ('pending', 'ready')",
		publicID,
	)
	if err != nil {
		slog.Error("failed to mark video ready", "error", err, "public_id", publicID)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("no pending or ready video found with id " + publicID.String())
	}

	slog.Info("marked video ready", "public_id", publicID)
	return nil
}

func (d *Database) CreateUser(name string, email string, password string) (int64, error) {
	var exists bool
	err := d.conn.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		slog.Error("error checking email", "error", err)
		return 0, err
	}

	if exists {
		return 0, errors.New("email already registered")
	}

	res := d.conn.QueryRow(
		"INSERT INTO users (name, email, hash) VALUES ($1, $2, $3) RETURNING id",
		name,
		email,
		password,
	)

	var id int64
	if err := res.Scan(&id); err != nil {
		slog.Error("failed to insert user", "error", err)
		return 0, err
	}

	slog.Info("Successfully added user:", name, email)
	return id, nil
}

func isValidYouTubeURL(url string) bool {
	return len(url) > 0 && ((len(url) > 23 && url[:24] == "https://www.youtube.com/") ||
		(len(url) > 19 && url[:20] == "http://www.youtube.com/") ||
		(len(url) > 16 && url[:17] == "https://youtu.be/") ||
		(len(url) > 15 && url[:16] == "http://youtu.be/"))
}

func extractYouTubeID(url string) string {
	if len(url) > 17 && url[:17] == "https://youtu.be/" {
		videoID := url[17:]
		if idx := indexOf(videoID, '?'); idx != -1 {
			videoID = videoID[:idx]
		}
		return videoID
	}
	if len(url) > 16 && url[:16] == "http://youtu.be/" {
		videoID := url[16:]
		if idx := indexOf(videoID, '?'); idx != -1 {
			videoID = videoID[:idx]
		}
		return videoID
	}

	vIndex := indexOfStr(url, "v=")
	if vIndex == -1 {
		return ""
	}

	videoID := url[vIndex+2:]
	if idx := indexOf(videoID, '&'); idx != -1 {
		videoID = videoID[:idx]
	}

	return videoID
}

func indexOf(s string, char byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == char {
			return i
		}
	}
	return -1
}

func indexOfStr(s string, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func (d *Database) Close() {
	d.conn.Close()
}
