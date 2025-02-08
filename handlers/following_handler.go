package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// func GetFollowingListHandler(db *database.Queries, user database.User) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		following, err := db.GetFollowingList(r.Context(), user.UserID)
// 		if err != nil {
// 			http.Error(w, "Couldn't get following list", http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK) // 200
// 		json.NewEncoder(w).Encode(following)
// 	})
// }

// func GetFollowerListHandler(db *database.Queries, user database.User) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		followers, err := db.GetFollowerList(r.Context(), user.UserID)
// 		if err != nil {
// 			http.Error(w, "Couldn't get follower list", http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK) // 200
// 		json.NewEncoder(w).Encode(followers)
// 	})
// }

func GetFollowingListByIDHandler(db *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		usernameParam := chi.URLParam(r, "username")
		username := usernameParam

		userID, err := db.GetIDbyUsername(r.Context(), username)
		if err != nil {
			http.Error(w, "Couldn't get user ID", http.StatusInternalServerError)
			return
		}

		following, err := db.GetFollowingList(r.Context(), userID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Couldn't get following list", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // 200
		json.NewEncoder(w).Encode(following)
	})
}

func GetFollowerListByIDHandler(db *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract userID from URL
		userIDParam := chi.URLParam(r, "userID")
		userID, err := uuid.Parse(userIDParam)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Get follower list for the user
		followers, err := db.GetFollowerList(r.Context(), userID)
		if err != nil {
			http.Error(w, "Couldn't get follower list", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // 200
		json.NewEncoder(w).Encode(followers)
	})
}

func CreateFollowHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			FollowingID uuid.UUID `json:"followee_id"`
		}

		var params parameters
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest) // 400
			return
		}

		err := db.CreateFollow(r.Context(), database.CreateFollowParams{
			FollowerID:  user.UserID,
			FollowingID: params.FollowingID,
		})
		if err != nil {
			http.Error(w, "Couldn't follow user", http.StatusInternalServerError) // 500
			return
		}

		w.WriteHeader(http.StatusCreated) // 201
	})
}

func DeleteFollowHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		FollowingID := chi.URLParam(r, "id")

		followeeUUID, err := uuid.Parse(FollowingID)
		if err != nil {
			http.Error(w, "Invalid followee ID", http.StatusBadRequest) // 400
			return
		}

		err = db.DeleteFollow(r.Context(), database.DeleteFollowParams{
			FollowerID:  user.UserID,
			FollowingID: followeeUUID,
		})
		if err != nil {
			http.Error(w, "Couldn't unfollow user", http.StatusInternalServerError) // 500
			return
		}

		w.WriteHeader(http.StatusOK) // 200
	})
}

// GetFeedHandler for fetching the feed
func GetFeedHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed, err := db.GetFeed(r.Context(), user.UserID)
		if err != nil {
			http.Error(w, "Couldn't get feed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // 200
		json.NewEncoder(w).Encode(feed)
	})
}
