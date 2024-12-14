package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/golang-jwt/jwt/v5"
)

func RefreshTokenHandler(db *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshCookie, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "Missing refresh token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(refreshCookie.Value, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			http.Error(w, "Invalid claims in refresh token", http.StatusUnauthorized)
			return
		}

		if time.Now().After(claims.ExpiresAt.Time) {
			http.Error(w, "Refresh token expired", http.StatusUnauthorized)
			return
		}

		accessToken, err := generateAccessToken(claims.UserID, claims.Role)
		if err != nil {
			http.Error(w, "Couldn't generate new access token", http.StatusInternalServerError)
			return
		}

		refreshToken, err := generateRefreshToken(claims.UserID)
		if err != nil {
			http.Error(w, "Couldn't generate new refresh token", http.StatusInternalServerError)
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
			"access_token": accessToken,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}
