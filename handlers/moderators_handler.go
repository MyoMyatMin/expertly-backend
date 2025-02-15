package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type returnedModerator struct {
	ModeratorID uuid.UUID `json:"moderator_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
}

func CreateModeratorHandler(db *database.Queries, moderator database.Moderator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Roles    string `json:"roles"`
		}

		var params parameters

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)

		if err != nil {
			http.Error(w, "Counldn't hash password", http.StatusInternalServerError)
		}

		moderator, err := db.CreateModerator(r.Context(), database.CreateModeratorParams{
			ModeratorID: uuid.New(),
			Name:        params.Name,
			Email:       params.Email,
			Password:    string(hashedPassword),
			CreatedBy:   uuid.NullUUID{UUID: moderator.ModeratorID, Valid: true},
			Role:        params.Roles,
		})

		if err != nil {
			http.Error(w, "Couldn't create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(moderator)
	})
}

func LoginModeratorController(db *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var params parameters

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		fmt.Println(params)

		moderator, err := db.GetModeratorByEmail(r.Context(), params.Email)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(moderator.Password), []byte(params.Password)); err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		accessToken, err := generateAccessToken(moderator.ModeratorID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Couldn't generate access token", http.StatusInternalServerError)
			return
		}

		refreshToken, err := generateRefreshToken(moderator.ModeratorID)
		if err != nil {
			http.Error(w, "Couldn't generate refresh token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Expires:  time.Now().Add(2 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})

		var returnedModerator = returnedModerator{
			ModeratorID: moderator.ModeratorID,
			Name:        moderator.Name,
			Email:       moderator.Email,
			Role:        moderator.Role,
		}

		response := map[string]interface{}{
			"user":          returnedModerator,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	})
}
