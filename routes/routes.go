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
)

func SetUpRoutes(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	godotenv.Load(".env")

	// CORS Middleware
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

	// Initialize database queries
	queries := database.New(db)

	// Public Routes
	r.Get("/", handlers.HelloHandler)
	r.Post("/signup", handlers.SignUpHandler(queries).ServeHTTP)
	r.Post("/login", handlers.LoginHandler(queries).ServeHTTP)
	r.Post("/logout", handlers.LogoutHandler)
	r.Post("/refresh_token", handlers.RefreshTokenHandler(queries).ServeHTTP)

	// User-specific Routes
	r.Get("/auth/me", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.CheckAuthStatsHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Get("/posts", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.GetAllPostsHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Put("/posts/{id}", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.UpdatePostHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Get("/posts/{id}", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.GetPostByIDHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Post("/posts/{postID}/upvotes", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.InsertUpvoteHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Delete("/posts/{postID}/upvotes", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.DeleteUpvoteHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Post("/posts/{postID}/comments", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.CreateCommentHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Delete("/posts/{postID}/comments/{commentID}", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.DeleteCommentHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Patch("/posts/{postID}/comments/{commentID}", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.UpdateCommentHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Get("/posts/{postID}/comments", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.GetAllCommentsByPostHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Get("/following", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.GetFollowingListHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Get("/followers", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.GetFollowerListHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Post("/follow", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.CreateFollowHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Delete("/follow/{id}", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.DeleteFollowHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Get("/feed", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.GetFeedHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Get("/users/{userID}/following", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.GetFollowingListByIDHandler(queries).ServeHTTP(w, r)
		}, nil, nil, "user"))

	r.Get("/users/{userID}/followers", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.GetFollowerListByIDHandler(queries).ServeHTTP(w, r)
		}, nil, nil, "user"))

	// Contributor Routes
	r.Post("/posts", middlewares.MiddlewareAuth(queries,
		nil,
		func(w http.ResponseWriter, r *http.Request, contributor database.Contributor) {
			handlers.CreatePostHandler(queries, contributor).ServeHTTP(w, r)
		}, nil, "contributor"))

	// Moderator Routes
	r.Post("/admin/create", middlewares.MiddlewareAuth(queries,
		nil, nil,
		func(w http.ResponseWriter, r *http.Request, moderator database.Moderator) {
			handlers.CreateModeratorHandler(queries, moderator).ServeHTTP(w, r)
		}, "moderator"))

	r.Post("/admin/login", handlers.LoginModeratorController(queries).ServeHTTP)

	return r
}
