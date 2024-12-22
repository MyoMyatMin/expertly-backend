package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func CreateCommentHandler(db *database.Queries, user database.User) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Define the parameters struct to receive JSON body
		type parameters struct {
			Content         string        `json:"content"`
			ParentCommentID uuid.NullUUID `json:"parent_comment_id"`
		}

		// Get the PostID from the URL path
		PostID := chi.URLParam(r, "post_id")

		// Parse the PostID to uuid.UUID
		postID, err := uuid.Parse(PostID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Decode the JSON body into the 'parameters' struct
		var params parameters
		err = json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Call CreateComment function from the database
		comment, err := db.CreateComment(r.Context(), database.CreateCommentParams{
			Content:         params.Content,
			ParentCommentID: params.ParentCommentID,
			PostID:          postID,
			UserID:          user.ID,
		})

		if err != nil {
			http.Error(w, "Failed to create comment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(comment); err != nil {
			http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		}
	})

}

func UpdateCommentHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Content string `json:"content"`
		}

		CommentID := chi.URLParam(r, "comment_id")
		commentID, err := uuid.Parse(CommentID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var params parameters
		err = json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		comment, err := db.UpdateComment(r.Context(), database.UpdateCommentParams{
			ID:      commentID,
			Content: params.Content,
		})
		if err != nil {
			http.Error(w, "Failed to update comment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(comment); err != nil {
			http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		}
	})

}

func DeleteCommentHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CommentID := chi.URLParam(r, "comment_id")
		commentID, err := uuid.Parse(CommentID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = db.DeleteComment(r.Context(), commentID)
		if err != nil {
			http.Error(w, "Failed to delete comment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
