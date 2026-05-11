package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gultekinmakif/go-http-server/internal/db/postgres"
	"github.com/gultekinmakif/go-http-server/internal/models"
	"github.com/gultekinmakif/go-http-server/internal/utils"
	"gorm.io/gorm"
)

// createPostRequest is the JSON body for POST /posts.
type createPostRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// CreatePost handles POST /posts.
// Body: {"title":"...","body":"..."}.
// 201 with the created post on success, 400 on bad body, 500 on db error.
func CreatePost(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	var body createPostRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if models.GenerateSlug(body.Title) == "" {
		http.Error(w, models.ErrSlugEmpty.Error(), http.StatusBadRequest)
		return
	}

	db := postgres.Get().WithContext(r.Context())
	slug, err := models.PickAvailableSlug(db, body.Title, uuid.Nil)
	if err != nil {
		slog.Error("pick slug failed", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	p := models.Post{Title: body.Title, Body: body.Body, Slug: slug}
	if err := db.Create(&p).Error; err != nil {
		slog.Error("post creation failed", "err", err)
		http.Error(w, "post creation failed", http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, p.Content())
}

// GetPost handles GET /posts/{slug}.
// 200 + post on hit, 404 on miss, 500 on db error.
func GetPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	var p models.Post
	err := postgres.Get().WithContext(r.Context()).First(&p, "slug = ?", slug).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("get post failed", "err", err, "slug", slug)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, p.Content())
}

// GetAllPost handles GET /posts.
// 200 + slice of all posts on success, 500 on db error.
func GetAllPost(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	if err := postgres.Get().WithContext(r.Context()).Find(&posts).Error; err != nil {
		slog.Error("list posts failed", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, utils.Map(posts, models.Post.Content))
}
