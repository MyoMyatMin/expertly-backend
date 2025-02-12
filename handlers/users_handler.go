package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/MyoMyatMin/expertly-backend/utils"
	"github.com/go-chi/chi/v5"
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

type ReturnedProfileUser struct {
	UserID         uuid.UUID `json:"user_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	SuspendedUntil time.Time `json:"suspended_until"`
	Role           string    `json:"role"`
	Followers      int       `json:"followers"`
	Following      int       `json:"following"`
	IsFollowing    bool      `json:"is_following"`
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

		returnUser := ReturnedUser{
			UserID:         user.UserID,
			Name:           user.Name,
			Email:          user.Email,
			Username:       user.Username,
			SuspendedUntil: user.SuspendedUntil.Time,
			Role:           "user",
		}

		response := map[string]interface{}{
			"user":          returnUser,
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
func CheckAuthStatsHandler(db *database.Queries, user database.User, moderator database.Moderator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var returnedUser ReturnedUser

		if user.UserID != uuid.Nil {
			isContributor, err := db.CheckIfUserIsContributor(r.Context(), user.UserID)
			if err != nil {
				http.Error(w, "Couldn't check if user is contributor", http.StatusInternalServerError)
				return
			}

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
		} else if moderator.ModeratorID != uuid.Nil {
			returnedUser = ReturnedUser{
				UserID: moderator.ModeratorID, // Assuming moderator has UserID
				Name:   moderator.Name,        // Assuming moderator has Name
				Email:  moderator.Email,       // Add logic if needed
				Role:   moderator.Role,        // Assuming moderator has Role
			}
		} else {
			http.Error(w, "No valid user or moderator provided", http.StatusBadRequest)
			return
		}

		response := map[string]interface{}{
			"user": returnedUser,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetProfileDataHandler(db *database.Queries, user database.User, moderator database.Moderator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		username := chi.URLParam(r, "username")

		aimedUser, err := db.GetUserByUsername(r.Context(), username)

		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		isContributor, err := db.CheckIfUserIsContributor(r.Context(), aimedUser.UserID)

		if err != nil {
			http.Error(w, "Couldn't check if user is contributor", http.StatusInternalServerError)
			return
		}

		follower_count, err := db.GetFollwersCount(r.Context(), aimedUser.UserID)
		if err != nil {
			http.Error(w, "Couldn't get follower count", http.StatusInternalServerError)
			return

		}

		following_count, err := db.GetFollowingCount(r.Context(), aimedUser.UserID)
		if err != nil {
			http.Error(w, "Couldn't get following count", http.StatusInternalServerError)
			return
		}

		isFollowing, err := db.GetFollowStatus(r.Context(), database.GetFollowStatusParams{
			FollowerID:  user.UserID,
			FollowingID: aimedUser.UserID,
		})

		if err != nil {
			http.Error(w, "Couldn't get follow status", http.StatusInternalServerError)
			return
		}

		returnedUser := ReturnedProfileUser{
			UserID:         aimedUser.UserID,
			Name:           aimedUser.Name,
			Email:          aimedUser.Email,
			Username:       aimedUser.Username,
			SuspendedUntil: aimedUser.SuspendedUntil.Time,
			Role:           "user",
			Followers:      int(follower_count),
			Following:      int(following_count),
			IsFollowing:    false,
		}
		if isFollowing {
			returnedUser.IsFollowing = true
		}

		if isContributor {
			returnedUser.Role = "contributor"
		}

		// followingList, err := db.GetFollowingList(r.Context(), aimedUser.UserID)
		// if err != nil {
		// 	http.Error(w, "Couldn't get following list", http.StatusInternalServerError)
		// 	return

		// }

		response := map[string]interface{}{
			"user": returnedUser,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetContributorProfilePostsHandler(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		aimedUser, err := db.GetUserByUsername(r.Context(), username)

		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		posts, err := db.GetPostsByContributor(r.Context(), aimedUser.UserID)
		if err != nil {
			http.Error(w, "Couldn't get posts by contributor", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

func UpdateUserHandler(db *database.Queries, user database.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Name     string `json:"name"`
			Username string `json:"username"`
		}

		var params parameters

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := db.UpdateUser(r.Context(), database.UpdateUserParams{
			UserID:   user.UserID,
			Name:     params.Name,
			Username: params.Username,
		})

		if err != nil {
			http.Error(w, "Couldn't update user", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"message": "User updated",
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
