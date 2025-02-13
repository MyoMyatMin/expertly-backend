package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func CreateReportHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Reason          string        `json:"reason"`
			TargetPostID    uuid.NullUUID `json:"target_postID"`
			TargetCommentID uuid.NullUUID `json:"target_CommentID"`
		}

		var params parameters

		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var targetUserID uuid.UUID
		var targetPostID uuid.NullUUID
		if params.TargetPostID.Valid {
			post, err := db.GetPost(r.Context(), params.TargetPostID.UUID)
			if err != nil {
				http.Error(w, "Post not found", http.StatusNotFound)
				return
			}
			targetUserID = post.UserID // Author of the post
			targetPostID = params.TargetPostID
		} else if params.TargetCommentID.Valid {
			comment, err := db.GetCommentByID(r.Context(), params.TargetCommentID.UUID)
			if err != nil {
				http.Error(w, "Comment not found", http.StatusNotFound)
				return
			}

			targetUserID = comment.UserID
			targetPostID = uuid.NullUUID{UUID: comment.PostID, Valid: true}

		} else {
			http.Error(w, "Invalid report target", http.StatusBadRequest)
			return
		}

		_, err = db.CreateReport(r.Context(), database.CreateReportParams{
			ReportID:        uuid.New(),
			ReportedBy:      user.UserID,
			TargetPostID:    targetPostID,
			TargetCommentID: params.TargetCommentID,
			TargetUserID:    targetUserID,
			Reason:          params.Reason,
		})

		if err != nil {
			http.Error(w, "Couldn't create report", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

func GetReportsHandler(db *database.Queries, moderator database.Moderator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetReportsHandler")
		reports, err := db.ListAllReportDetails(r.Context())
		if err != nil {
			fmt.Print(err)
			http.Error(w, "Couldn't get reports", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(reports)
	})
}

func UpdateReportStatusHandler(db *database.Queries, moderator database.Moderator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Status string `json:"status"`
		}

		reportIDStr := chi.URLParam(r, "reportID")
		reportID, err := uuid.Parse(reportIDStr)
		if err != nil {
			http.Error(w, "Invalid report ID", http.StatusBadRequest)
			return
		}

		var params parameters
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		fmt.Println("UpdateReportStatusHandler", params, reportID)
		if !validStatuses[params.Status] {
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

		status := sql.NullString{String: params.Status, Valid: params.Status != ""}
		_, err = db.UpdateReportStatus(r.Context(), database.UpdateReportStatusParams{
			ReportID:   reportID,
			Status:     status,
			Reviewedby: uuid.NullUUID{UUID: moderator.ModeratorID, Valid: true},
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Couldn't update report status", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Report status updated successfully"})
	})
}

func GetReportedContributorsHandler(db *database.Queries, moderator database.Moderator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetReportedContributorsHandler")
		reports, err := db.ListReportedContributors(r.Context())
		if err != nil {
			fmt.Print(err)
			http.Error(w, "Couldn't get reports", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(reports)
	})
}

func GetReportedUserHandler(db *database.Queries, moderator database.Moderator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reports, err := db.ListReportedUsers(r.Context())
		if err != nil {
			fmt.Print(err)
			http.Error(w, "Couldn't get reports", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(reports)
	})
}
