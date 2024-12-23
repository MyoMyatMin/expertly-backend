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
	queries := database.New(db)
	r.Get("/", handlers.HelloHandler)
	r.Post("/signup", handlers.SignUpHandler(queries).ServeHTTP)
	r.Post("/login", handlers.LoginHandler(queries).ServeHTTP)
	r.Post("/logout", handlers.LogoutHandler)
	r.Post("/refresh_token", handlers.RefreshTokenHandler(queries).ServeHTTP)
	r.Get("/auth/me", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.CheckAuthStatsHandler(queries, user).ServeHTTP(w, r)
		}))
	r.Get("/test_middlewares", middlewares.MiddlewareAuth(queries, func(w http.ResponseWriter, r *http.Request, user database.User) {
		handlers.TestMiddlewaresHandler(queries, user).ServeHTTP(w, r)
	}))

	r.Get("/posts", middlewares.MiddlewareAuth(queries, func(w http.ResponseWriter, r *http.Request, user database.User) {
		handlers.GetAllPostsHandler(queries, user).ServeHTTP(w, r)
	}))

	r.Post("/posts", middlewares.MiddlewareAuth(queries, func(w http.ResponseWriter, r *http.Request, user database.User) {
		handlers.CreatePostHandler(queries, user).ServeHTTP(w, r)
	}))

	r.Put("/posts/{id}", middlewares.MiddlewareAuth(queries, func(w http.ResponseWriter, r *http.Request, user database.User) {
		handlers.UpdatePostHandler(queries, user).ServeHTTP(w, r)
	}))

	r.Get("/posts/{id}", middlewares.MiddlewareAuth(queries, func(w http.ResponseWriter, r *http.Request, user database.User) {
		handlers.GetPostByIDHandler(queries, user).ServeHTTP(w, r)
	}))

	r.Post("/posts/{postID}/upvotes", middlewares.MiddlewareAuth(queries, func(w http.ResponseWriter, r *http.Request, user database.User) {
		handlers.InsertUpvoteHandler(queries, user).ServeHTTP(w, r)
	}))

	r.Delete("/posts/{postID}/upvotes", middlewares.MiddlewareAuth(queries, func(w http.ResponseWriter, r *http.Request, user database.User) {
		handlers.DeleteUpvoteHandler(queries, user).ServeHTTP(w, r)
	}))

	r.Post("/posts/{postID}/comments", middlewares.MiddlewareAuth(queries, func(w http.ResponseWriter, r *http.Request, user database.User) {
		handlers.CreateCommentHandler(queries, user).ServeHTTP(w, r)
	}))

	r.Delete("/posts/{postID}/comments/{commentID}", middlewares.MiddlewareAuth(queries, func(w http.ResponseWriter, r *http.Request, user database.User) {
		handlers.DeleteCommentHandler(queries, user).ServeHTTP(w, r)
	}))

	r.Patch("/posts/{postID}/comments/{commentID}", middlewares.MiddlewareAuth(queries, func(w http.ResponseWriter, r *http.Request, user database.User) {
		handlers.UpdateCommentHandler(queries, user).ServeHTTP(w, r)
	}))

	r.Get("/posts/{postID}/comments", middlewares.MiddlewareAuth(queries, func(w http.ResponseWriter, r *http.Request, user database.User) {
		handlers.GetAllCommentsByPostHandler(queries, user).ServeHTTP(w, r)
	}))

	return r
}
