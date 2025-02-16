package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

		fmt.Println(params)

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
			fmt.Println(err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		appealIDStr := chi.URLParam(r, "appealID")
		appealID, err := uuid.Parse(appealIDStr)

		if err != nil {
			fmt.Println(err, appealIDStr)
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

		if params.Status == "resolved" {
			appeal, err := db.GetAppealById(r.Context(), appealID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			report, err := db.GetReportById(r.Context(), appeal.TargetReportID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			user, err := db.GetUserById(r.Context(), appeal.AppealedBy)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var newSuspendedUntil sql.NullTime
			if report.SuspendDays.Int32 == 0 {
				newSuspendedUntil = sql.NullTime{Valid: false}
			} else {
				newSuspendedUntil = sql.NullTime{Time: user.SuspendedUntil.Time.AddDate(0, 0, -int(report.SuspendDays.Int32)), Valid: true}
				if newSuspendedUntil.Time.Before(time.Now()) {
					newSuspendedUntil = sql.NullTime{Valid: false}
				}
			}
			err = db.UpdateUserSuspension(r.Context(), database.UpdateUserSuspensionParams{
				UserID:         user.UserID,
				SuspendedUntil: newSuspendedUntil,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
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

func GetContributorsAppeals(db *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		appeals, err := db.ListAppealsByContributors(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(appeals)
	})
}

func GetUsersAppeals(db *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		appeals, err := db.ListAppealsByUsers(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(appeals)
	})
}
