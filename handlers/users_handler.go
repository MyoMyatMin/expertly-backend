package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = os.Getenv("SECRET_KEY")

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func generateAccessToken(userID uuid.UUID, role string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

func generateRefreshToken(userID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

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

		accessToken, err := generateAccessToken(user.ID, role)
		if err != nil {
			http.Error(w, "Couldn't generate access token", http.StatusInternalServerError)
			return
		}

		refreshToken, err := generateRefreshToken(user.ID)
		if err != nil {
			http.Error(w, "Couldn't generate refresh token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Expires:  time.Now().Add(1 * time.Hour),
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

		response := map[string]interface{}{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
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

		accessToken, err := generateAccessToken(user.ID, user.Role.String)
		if err != nil {
			http.Error(w, "Couldn't generate access token", http.StatusInternalServerError)
			return
		}

		refreshToken, err := generateRefreshToken(user.ID)
		if err != nil {
			http.Error(w, "Couldn't generate refresh token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Expires:  time.Now().Add(1 * time.Hour),
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

		response := map[string]interface{}{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}
