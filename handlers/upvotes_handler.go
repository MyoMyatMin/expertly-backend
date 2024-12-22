package handlers

import (
	"net/http"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func InsertUpvoteHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "postID")

		postUUID, err := uuid.Parse(postID)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		upvote := database.Upvote{
			UserID: user.ID,
			PostID: postUUID,
		}

		_, err = db.InsertUpvote(r.Context(), database.InsertUpvoteParams{
			UserID: upvote.UserID,
			PostID: upvote.PostID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	})
}

func DeleteUpvoteHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "postID")

		postUUID, err := uuid.Parse(postID)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		upvote := database.Upvote{
			UserID: user.ID,
			PostID: postUUID,
		}

		_, err = db.DeleteUpvote(r.Context(), database.DeleteUpvoteParams{
			UserID: upvote.UserID,
			PostID: upvote.PostID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
