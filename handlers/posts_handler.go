package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func GetAllPostsHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		posts, err := db.ListPosts(r.Context())
		if err != nil {
			http.Error(w, "Couldn't get posts", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(posts)
	})
}

func CreatePostHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Title   string `json:"title"`
			Content string `json:"content"`
			Status  string `json:"status"`
		}

		var params parameters

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		validStatuses := map[string]bool{
			"draft":     true,
			"published": true,
		}

		if _, valid := validStatuses[params.Status]; !valid {
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

		statusNull := sql.NullString{
			String: params.Status,
			Valid:  true,
		}

		userIDNull := uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		}

		post, err := db.CreatePost(r.Context(), database.CreatePostParams{
			ID:      uuid.New(),
			Title:   params.Title,
			Content: params.Content,
			Status:  statusNull,
			UserID:  userIDNull,
		})
		if err != nil {
			http.Error(w, "Couldn't create post", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(post)
	})

}

func GetPostByIDHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		postID := chi.URLParam(r, "id")

		postUUID, err := uuid.Parse(postID)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		post, err := db.GetPost(r.Context(), postUUID)
		if err != nil {
			http.Error(w, "Couldn't get post", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(post)
	})
}

func UpdatePostHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "id")

		postUUID, err := uuid.Parse(postID)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		type parameters struct {
			Title   string `json:"title"`
			Content string `json:"content"`
			Status  string `json:"status"`
		}

		var params parameters

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		validStatuses := map[string]bool{
			"draft":     true,
			"published": true,
		}

		if _, valid := validStatuses[params.Status]; !valid {
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

		statusNull := sql.NullString{
			String: params.Status,
			Valid:  true,
		}

		post, err := db.UpdatePost(r.Context(), database.UpdatePostParams{
			ID:      postUUID,
			Title:   params.Title,
			Content: params.Content,
			Status:  statusNull,
		})
		if err != nil {
			http.Error(w, "Couldn't update post", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(post)
	})
}
