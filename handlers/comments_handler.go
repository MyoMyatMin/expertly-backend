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
	Username        string        `json:"username"`
	Name            string        `json:"name"`
	Replies         []*Comment    `json:"replies"` // Change to slice of pointers
}

func CreateCommentHandler(db *database.Queries, user database.User) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type parameters struct {
			Content         string        `json:"content"`
			ParentCommentID uuid.NullUUID `json:"replyCommentId"`
		}
		PostID := chi.URLParam(r, "postID")

		postID, err := uuid.Parse(PostID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var params parameters
		fmt.Println(r.Body, params)
		err = json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			fmt.Println(err)
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
		fmt.Println("Here")
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

func GetAllCommentsByPostHandler(db *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postSlug := chi.URLParam(r, "postSlug")
		PostID, err := db.GetPostBySlug(r.Context(), postSlug)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to get post ID: "+err.Error(), http.StatusInternalServerError)
			return
		}

		postID := PostID.PostID

		dbcomments, err := db.GetCommentsByPost(r.Context(), postID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to get comments: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(len(dbcomments))

		var comments []Comment

		for _, dbcomment := range dbcomments {
			comments = append(comments, Comment{
				ID:              dbcomment.CommentID,
				Content:         dbcomment.Content,
				ParentCommentID: dbcomment.ParentCommentID,
				PostID:          dbcomment.PostID,
				UserID:          dbcomment.UserID,
				Username:        dbcomment.Username,
				Name:            dbcomment.Name,
			})

		}

		nestedComments := BuildNestedComments(comments)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(nestedComments); err != nil {
			http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		}
	})
}

func BuildNestedComments(comments []Comment) []*Comment {
	commentMap := make(map[uuid.UUID]*Comment)

	// Initialize map with pointers to each comment in the original slice
	for i := range comments {
		comment := &comments[i]        // Pointer to the original comment in the slice
		comment.Replies = []*Comment{} // Initialize Replies as slice of pointers
		commentMap[comment.ID] = comment
	}

	var nestedComments []*Comment

	for i := range comments {
		comment := commentMap[comments[i].ID]
		if !comment.ParentCommentID.Valid {
			// Add top-level comments to the result
			nestedComments = append(nestedComments, comment)
		} else {
			// Append the comment as a reply to its parent
			parentID := comment.ParentCommentID.UUID
			if parentComment, ok := commentMap[parentID]; ok {
				parentComment.Replies = append(parentComment.Replies, comment)
			}
		}
	}

	return nestedComments
}
