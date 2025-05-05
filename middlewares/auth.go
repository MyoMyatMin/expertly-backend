package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type HandlerWithUser func(http.ResponseWriter, *http.Request, database.User)
type HandlerWithContributor func(http.ResponseWriter, *http.Request, database.Contributor)
type HandlerWithModerator func(http.ResponseWriter, *http.Request, database.Moderator)

func loadEnvIfLocal() {
	if os.Getenv("Local") == "local" {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("Warning: failed to load .env (this is okay in production)")
		}
	}
}

// MiddlewareAuth handles role-based auth
func MiddlewareAuth(
	db *database.Queries,
	handlerWithUser HandlerWithUser,
	handlerWithContributor HandlerWithContributor,
	handlerWithModerator HandlerWithModerator,
	authType string, // "user", "contributor", or "moderator"
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loadEnvIfLocal()

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

		switch authType {
		case "moderator":
			moderatorRow, err := db.GetModeratorById(r.Context(), userID)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Moderator not found")
				return
			}
			moderator := database.Moderator{
				ModeratorID: moderatorRow.ModeratorID,
				CreatedAt:   moderatorRow.CreatedAt,
				Role:        moderatorRow.Role,
				Email:       moderatorRow.Email,
				Name:        moderatorRow.Name,
			}
			handlerWithModerator(w, r, moderator)

		case "contributor":
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
			handlerWithContributor(w, r, contributor)

		case "user":
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
			handlerWithUser(w, r, user)

		default:
			respondWithError(w, http.StatusBadRequest, "Invalid authentication type")
		}
	}
}

// MiddlewareModeratorOrUser allows both roles
func MiddlewareModeratorOrUser(
	db *database.Queries,
	handlerWithUser HandlerWithUser,
	handlerWithModerator HandlerWithModerator,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loadEnvIfLocal()

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

		// Try moderator first
		if moderatorRow, err := db.GetModeratorById(r.Context(), userID); err == nil {
			moderator := database.Moderator{
				ModeratorID: moderatorRow.ModeratorID,
				Role:        moderatorRow.Role,
				Email:       moderatorRow.Email,
				Name:        moderatorRow.Name,
				CreatedAt:   moderatorRow.CreatedAt,
			}
			handlerWithModerator(w, r, moderator)
			return
		}

		// Fall back to user
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
		handlerWithUser(w, r, user)
	}
}

// --- Utility functions ---

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
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
	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		return nil, errors.New("missing JWT secret key")
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if isTokenExpired(claims) {
		return nil, errors.New("token is expired")
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
	exp, ok := claims["exp"].(float64)
	if !ok {
		return true
	}
	return time.Now().Unix() > int64(exp)
}
