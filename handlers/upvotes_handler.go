package handlers

import (
	"net/http"

	"encoding/json"

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
			UserID: user.UserID,
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
			UserID: user.UserID,
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

func GetUpvotesHandlerByUser(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		upvotes, err := db.ListUpvotesByUser(r.Context(), user.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(upvotes)
	})

}
