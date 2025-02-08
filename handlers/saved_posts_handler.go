package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func CreateSavePost(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			PostID string `json:"post_id"`
		}

		var params parameters

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest) // 400
			return
		}

		postUUID, err := uuid.Parse(params.PostID)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest) // 400
			return
		}

		_, err = db.CreateSavedPost(r.Context(), database.CreateSavedPostParams{
			UserID: user.UserID,
			PostID: postUUID,
		})
		if err != nil {
			http.Error(w, "Couldn't save post", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // 200
	})
}

func DeleteSavedPost(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "id")

		postUUID, err := uuid.Parse(postID)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest) // 400
			return
		}

		err = db.DeleteSavedPost(r.Context(), database.DeleteSavedPostParams{
			UserID: user.UserID,
			PostID: postUUID,
		})
		if err != nil {
			http.Error(w, "Couldn't delete saved post", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // 200
	})
}

func GetSavedPosts(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		userID, err := db.GetIDbyUsername(r.Context(), username)
		if err != nil {
			http.Error(w, "Couldn't get user ID", http.StatusInternalServerError)
			return
		}

		savedPosts, err := db.ListSavedPostsByID(r.Context(), userID)
		if err != nil {
			http.Error(w, "Couldn't get saved posts", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // 200
		json.NewEncoder(w).Encode(savedPosts)

	})

}
