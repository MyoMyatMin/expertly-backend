package routes

import (
	"database/sql"
	"net/http"

	"github.com/MyoMyatMin/expertly-backend/handlers"
	"github.com/MyoMyatMin/expertly-backend/middlewares"
	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	// Fetch user from the database
)

func SetUpRoutes(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	godotenv.Load(".env")
	r.Use(
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
	)
	quaries := database.New(db)
	r.Get("/", handlers.HelloHandler)
	r.Post("/signup", handlers.SignUpHandler(quaries).ServeHTTP)
	r.Post("/login", handlers.LoginHandler(quaries).ServeHTTP)
	r.Post("/logout", handlers.LogoutHandler)
	r.Post("/refresh_token", handlers.RefreshTokenHandler(quaries).ServeHTTP)
	r.Get("/auth/me", middlewares.MiddlewareAuth(quaries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.CheckAuthStatsHandler(quaries, user).ServeHTTP(w, r)
		}))
	r.Get("/test_middlewares", middlewares.MiddlewareAuth(quaries, func(w http.ResponseWriter, r *http.Request, user database.User) {
		handlers.TestMiddlewaresHandler(quaries, user).ServeHTTP(w, r)
	}))

	return r
}
