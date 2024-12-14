package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SignUpHandler(db *database.Queries) http.Handler {
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

		validRoles := map[string]bool{
			"user":        true,
			"contributor": true,
			"moderator":   true,
			"superadmin":  true,
		}

		role := "user"

		if params.Roles != "" && validRoles[params.Roles] {
			role = params.Roles
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Couldn't hash password", http.StatusInternalServerError)
			return
		}

		roleNull := sql.NullString{String: role, Valid: true}

		user, err := db.CreateUser(r.Context(), database.CreateUserParams{
			ID:       uuid.New(),
			Name:     params.Name,
			Email:    params.Email,
			Password: string(passwordHash),
			Role:     roleNull,
		})

		if err != nil {
			http.Error(w, "Couldn't create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

	})

}

func LoginHandler(db *database.Queries) http.Handler {
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

		user, err := db.GetUserByEmail(r.Context(), params.Email)

		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))

		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})
}
