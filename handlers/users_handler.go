package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/MyoMyatMin/expertly-backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type ReturnedUser struct {
	UserID         uuid.UUID `json:"user_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	SuspendedUntil time.Time `json:"suspended_until"`
	Role           string    `json:"role"`
}

func generateAccessToken(userID uuid.UUID) (string, error) {
	godotenv.Load(".env")
	var jwtSecretKey = os.Getenv("SECRET_KEY")
	expirationTime := time.Now().Add(1 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"user_id": userID,

		"exp": expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecretKey))
}

func generateRefreshToken(userID uuid.UUID) (string, error) {
	godotenv.Load(".env")
	var jwtSecretKey = os.Getenv("SECRET_KEY")
	expirationTime := time.Now().Add(24 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"user_id": userID,

		"exp": expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecretKey))

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

		username, err := utils.GenerateUniqueUsername(params.Name, db, r)
		if err != nil {
			http.Error(w, "Couldn't generate unique username", http.StatusInternalServerError)
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Couldn't hash password", http.StatusInternalServerError)
			return
		}

		user, err := db.CreateUser(r.Context(), database.CreateUserParams{
			UserID:   uuid.New(),
			Name:     params.Name,
			Email:    params.Email,
			Password: string(passwordHash),
			Username: username,
		})

		if err != nil {
			http.Error(w, "Couldn't create user", http.StatusInternalServerError)
			return
		}

		accessToken, err := generateAccessToken(user.UserID)
		if err != nil {
			http.Error(w, "Couldn't generate access token", http.StatusInternalServerError)
			return
		}

		refreshToken, err := generateRefreshToken(user.UserID)
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

		accessToken, err := generateAccessToken(user.UserID)
		if err != nil {
			http.Error(w, "Couldn't generate access token", http.StatusInternalServerError)
			return
		}

		refreshToken, err := generateRefreshToken(user.UserID)
		if err != nil {
			http.Error(w, "Couldn't generate refresh token", http.StatusInternalServerError)
			return
		}

		var returnedUser ReturnedUser

		isContributor, err := db.CheckIfUserIsContributor(r.Context(), user.UserID)

		if err != nil {
			http.Error(w, "Couldn't check if user is contributor", http.StatusInternalServerError)
			return
		}

		returnedUser = ReturnedUser{
			UserID:         user.UserID,
			Name:           user.Name,
			Email:          user.Email,
			Username:       user.Username,
			SuspendedUntil: user.SuspendedUntil.Time,
			Role:           "user",
		}

		if isContributor {
			returnedUser.Role = "contributor"
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

		response := map[string]interface{}{
			"user":          returnedUser,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}

// func CheckAuthStatsHander(db *database.Queries, user database.User) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		accessToken, err := r.Cookie("access_token")
// 		if err != nil {
// 			http.Error(w, "No access token", http.StatusUnauthorized)
// 			return
// 		}

// 		claims := &JWTClaims{}
// 		_, err = jwt.ParseWithClaims(accessToken.Value, claims, func(token *jwt.Token) (interface{}, error) {
// 			return []byte(os.Getenv("SECRET_KEY")), nil
// 		})

// 		if err != nil {
// 			http.Error(w, "Invalid access token", http.StatusUnauthorized)
// 			return
// 		}

// 		user, err := db.GetUserById(r.Context(), claims.UserID)
// 		if err != nil {
// 			http.Error(w, "User not found", http.StatusNotFound)
// 			return
// 		}

// 		response := map[string]interface{}{
// 			"user": user,
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(response)
// 	})
// }

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	})

	response := map[string]interface{}{
		"message": "Logged out",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
func CheckAuthStatsHandler(db *database.Queries, user database.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		isContributor, err := db.CheckIfUserIsContributor(r.Context(), user.UserID)

		if err != nil {
			http.Error(w, "Couldn't check if user is contributor", http.StatusInternalServerError)
			return
		}

		var returnedUser ReturnedUser

		returnedUser = ReturnedUser{
			UserID:   user.UserID,
			Name:     user.Name,
			Email:    user.Email,
			Username: user.Username,
			Role:     "user",
		}

		if isContributor {
			returnedUser.Role = "contributor"
		}

		response := map[string]interface{}{
			"user": returnedUser,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// Modify TestMiddlewaresHandler to accept a user parameter
func TestMiddlewaresHandler(db *database.Queries, user database.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"user": user,
		}
		fmt.Println("User: ", user)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
