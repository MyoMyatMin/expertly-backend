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

// Define handler types for User and Contributor
type HandlerWithUser func(http.ResponseWriter, *http.Request, database.User)
type HandlerWithContributor func(http.ResponseWriter, *http.Request, database.Contributor)

// MiddlewareAuth is a generic middleware that works for both User and Contributor
func MiddlewareAuth(
	db *database.Queries,
	handlerWithUser HandlerWithUser,
	handlerWithContributor HandlerWithContributor,
	isContributor bool, // Flag to determine if the middleware is for Contributor
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract and validate the JWT token
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

		// Extract user ID from claims
		userID, err := getUserIDFromClaims(claims)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Fetch the user or contributor based on the flag
		if isContributor {
			// Fetch contributor
			contributorRow, err := db.GetContributorByUserId(r.Context(), userID)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Contributor not found")
				return
			}

			contributor := database.Contributor{
				UserID:          contributorRow.UserID,
				ExpertiseFields: contributorRow.ExpertiseFields,
				CreatedAt:       contributorRow.CreatedAt,
			}

			// Call the Contributor-specific handler
			handlerWithContributor(w, r, contributor)
		} else {
			// Fetch user
			userRow, err := db.GetUserById(r.Context(), userID)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "User not found")
				return
			}

			user := database.User{
				UserID:         userRow.UserID,
				Name:           userRow.Name,
				Username:       userRow.Username,
				Email:          userRow.Email,
				Password:       userRow.Password,
				SuspendedUntil: userRow.SuspendedUntil,
				CreatedAt:      userRow.CreatedAt,
				UpdatedAt:      userRow.UpdatedAt,
			}

			// Call the User-specific handler
			handlerWithUser(w, r, user)
		}
	}
}

// Helper functions remain unchanged
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

func isTokenExpired(claims jwt.MapClaims) bool {
	expTime, ok := claims["exp"].(float64)
	if !ok {
		return true // If there's no "exp" claim, consider it invalid
	}
	expTimeUnix := int64(expTime)
	currentTime := time.Now().Unix()
	return currentTime > expTimeUnix
}
