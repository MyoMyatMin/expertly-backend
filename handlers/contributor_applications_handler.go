package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/go-chi/chi/v5"

	"github.com/google/uuid"
)

func CreateContributorApplication(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			ExpertiseProofs   []string `json:"expertise_proofs"`
			IdentityProof     string   `json:"identity_proof"`
			InitialSubmission string   `json:"initial_submission"`
		}

		var params parameters

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		applicaition, err := db.ApplyContributorApplication(r.Context(), database.ApplyContributorApplicationParams{
			UserID:            user.UserID,
			ExpertiseProofs:   params.ExpertiseProofs,
			IdentityProof:     params.IdentityProof,
			InitialSubmission: params.InitialSubmission,
		})
		if err != nil {
			http.Error(w, "Couldn't create application", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated) // 201
		json.NewEncoder(w).Encode(applicaition)
	})
}

func UpdateContributorApplication(db *database.Queries, moderator database.Moderator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contri_app_id := chi.URLParam(r, "id")
		parsedID, err := uuid.Parse(contri_app_id)
		if err != nil {
			http.Error(w, "Invalid contributor application ID", http.StatusBadRequest)
			return
		}

		type parameters struct {
			Status string `json:"status"`
		}

		var params parameters

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		status := sql.NullString{String: params.Status, Valid: true}
		application, err := db.UpdateContributorApplication(r.Context(), database.UpdateContributorApplicationParams{
			ContriAppID: parsedID,
			Status:      status,
			ReviewedBy:  uuid.NullUUID{UUID: moderator.ModeratorID, Valid: true},
		})
		if err != nil {
			http.Error(w, "Couldn't update application", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // 200
		json.NewEncoder(w).Encode(application)

	})
}

func GetContributorApplications(db *database.Queries, moderator database.Moderator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		applications, err := db.ListContributorApplications(r.Context())
		if err != nil {
			http.Error(w, "Couldn't get applications", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // 200
		json.NewEncoder(w).Encode(applications)
	})
}

func GetContributorApplicationByID(db *database.Queries, moderator database.Moderator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contri_app_id := chi.URLParam(r, "id")
		parsedID, err := uuid.Parse(contri_app_id)
		if err != nil {
			http.Error(w, "Invalid contributor application ID", http.StatusBadRequest)
			return
		}

		application, err := db.GetContributorApplication(r.Context(), parsedID)
		if err != nil {
			http.Error(w, "Couldn't get application", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK) // 200
		json.NewEncoder(w).Encode(application)
	})
}
