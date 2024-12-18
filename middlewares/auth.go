package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type HandlerWithUser func(http.ResponseWriter, *http.Request, database.User)

func MiddlewareAuth(db *database.Queries, handler HandlerWithUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := extractTokenCookie(r)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		claims, err := parseJWTToken(tokenString)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		userID, err := getUserIDFromClaims(claims)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		user, err := db.GetUserById(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "User not found")
			return
		}

		handler(w, r, user)
	}
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]string{"error": message}
	_ = json.NewEncoder(w).Encode(response)
}

func extractTokenCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", errors.New("missing access token cookie")
		}
		return "", fmt.Errorf("could not retrieve access token cookie: %v", err)
	}
	if cookie.Value == "" {
		return "", errors.New("empty access token cookie")
	}
	return cookie.Value, nil
}

func parseJWTToken(tokenString string) (jwt.MapClaims, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("failed to load .env file")
	}
	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		return nil, errors.New("missing JWT secret key")
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if isTokenExpired(claims) {
		return nil, errors.New("token is expired")
	}

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func getUserIDFromClaims(claims jwt.MapClaims) (uuid.UUID, error) {
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("invalid token claims: user_id missing or not a string")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID format in token")
	}
	return userID, nil
}

// isTokenExpired checks if the JWT token has expired.
func isTokenExpired(claims jwt.MapClaims) bool {
	// Retrieve expiration time from claims
	expTime, ok := claims["exp"].(float64)
	if !ok {
		return true // If there's no "exp" claim, consider it invalid
	}

	// Compare expiration time with current time
	expTimeUnix := int64(expTime)
	currentTime := time.Now().Unix()

	return currentTime > expTimeUnix
}
