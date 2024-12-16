package routes

import (
	"database/sql"

	"github.com/MyoMyatMin/expertly-backend/handlers"
	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
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
	r.Post("/refresh_token", handlers.RefreshTokenHandler(quaries).ServeHTTP)
	r.Get("/auth/me", handlers.CheckAuthStatsHander(quaries).ServeHTTP)

	return r
}
