package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Comment struct {
	ID              uuid.UUID     `json:"id"`
	Content         string        `json:"content"`
	ParentCommentID uuid.NullUUID `json:"parent_comment_id"`
	PostID          uuid.UUID     `json:"post_id"`
	UserID          uuid.UUID     `json:"user_id"`
	Replies         []Comment     `json:"replies"`
}

func CreateCommentHandler(db *database.Queries, user database.User) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type parameters struct {
			Content         string        `json:"content"`
			ParentCommentID uuid.NullUUID `json:"parent_comment_id"`
		}
		PostID := chi.URLParam(r, "postID")

		postID, err := uuid.Parse(PostID)
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

		comment, err := db.CreateComment(r.Context(), database.CreateCommentParams{
			CommentID:       uuid.New(),
			Content:         params.Content,
			ParentCommentID: params.ParentCommentID,
			PostID:          postID,
			UserID:          user.UserID,
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

		CommentID := chi.URLParam(r, "commentID")
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
			CommentID: commentID,
			Content:   params.Content,
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
		CommentID := chi.URLParam(r, "commentID")
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

func GetAllCommentsByPostHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		PostID := chi.URLParam(r, "postID")
		postID, err := uuid.Parse(PostID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		dbcomments, err := db.GetCommentsByPost(r.Context(), postID)
		if err != nil {
			http.Error(w, "Failed to get comments: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var comments []Comment

		for _, dbcomment := range dbcomments {
			comments = append(comments, Comment{
				ID:              dbcomment.CommentID,
				Content:         dbcomment.Content,
				ParentCommentID: dbcomment.ParentCommentID,
				PostID:          dbcomment.PostID,
				UserID:          dbcomment.UserID,
			})

		}

		nestedComments := BuildNestedComments(comments)

		fmt.Println(nestedComments)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(nestedComments); err != nil {
			http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		}
	})
}

func BuildNestedComments(comments []Comment) []*Comment {
	commentMap := make(map[uuid.UUID]*Comment)

	for i := range comments {
		commentMap[comments[i].ID] = &comments[i]
	}

	var nestedComments []*Comment

	for i := range comments {
		comment := &comments[i]

		if !comment.ParentCommentID.Valid {
			nestedComments = append(nestedComments, comment)
		} else {
			parentID := comment.ParentCommentID.UUID
			parentComment, ok := commentMap[parentID]

			if ok {
				parentComment.Replies = append(parentComment.Replies, *comment)

			} else {
				fmt.Println("Parent comment not found")
			}
		}
	}

	return nestedComments
}
