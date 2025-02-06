package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func CreateAppealHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Reason         string    `json:"reason"`
			TargetReportID uuid.UUID `json:"target_reportID"`
		}

		var params parameters

		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		_, err = db.CreateAppeal(r.Context(), database.CreateAppealParams{
			AppealID:       uuid.New(),
			AppealedBy:     user.UserID,
			TargetReportID: params.TargetReportID,
			Reason:         params.Reason,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

var validStatuses = map[string]bool{
	"resolved":  true,
	"dismissed": true,
}

func UpdateAppealStatus(db *database.Queries, moderator database.Moderator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Status string `json:"status"`
		}

		var params parameters
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		appealIDStr := chi.URLParam(r, "appealID")
		appealID, err := uuid.Parse(appealIDStr)

		if err != nil {
			http.Error(w, "Invalid appeal ID", http.StatusBadRequest)
			return
		}

		status := sql.NullString{String: params.Status, Valid: params.Status != ""}
		_, err = db.UpdateAppealStatus(r.Context(), database.UpdateAppealStatusParams{
			AppealID:   appealID,
			Status:     status,
			Reviewedby: uuid.NullUUID(uuid.NullUUID{UUID: moderator.ModeratorID, Valid: true}),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func GetAppealsHandler(db *database.Queries, moderator database.Moderator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appeals, err := db.ListAllAppealDetails(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(appeals)
	})
}

func GetAppealByIDHandler(db *database.Queries, moderator database.Moderator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appealIDStr := chi.URLParam(r, "appealID")
		appealID, err := uuid.Parse(appealIDStr)
		if err != nil {
			http.Error(w, "Invalid appeal ID", http.StatusBadRequest)
			return
		}

		appeal, err := db.GetAppealById(r.Context(), appealID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(appeal)
	})
}
