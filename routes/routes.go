// routes.go
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

	r.Use(
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://expertly-psi.vercel.app", "http://localhost:3000"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders:   []string{"Content-Type", "Authorization"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
	)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Expertly Backend</title>
			</head>
			<body>
				<h1>Welcome to Expertly Backend</h1>
				<p>API documentation: <a href="https://github.com/MyoMyatMin/expertly-backend/blob/main/routes/routes.go" target="_blank">Routes Reference</a></p>
			</body>
			</html>
		`))
	})
	// Create API router
	apiRouter := chi.NewRouter()
	r.Mount("/api", apiRouter)

	queries := database.New(db)

	apiRouter.Post("/auth/signup", handlers.SignUpHandler(queries).ServeHTTP)
	apiRouter.Post("/auth/login", handlers.LoginHandler(queries).ServeHTTP)
	apiRouter.Post("/auth/logout", handlers.LogoutHandler)
	apiRouter.Post("/auth/refresh-token", handlers.RefreshTokenHandler(queries).ServeHTTP)
	apiRouter.Get("/auth/me", middlewares.MiddlewareModeratorOrUser(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.CheckAuthStatsHandler(queries, user, database.Moderator{}).ServeHTTP(w, r)
		},
		func(w http.ResponseWriter, r *http.Request, moderator database.Moderator) {
			handlers.CheckAuthStatsHandler(queries, database.User{}, moderator).ServeHTTP(w, r)
		}))

	// Post Routes
	apiRouter.Get("/posts", handlers.GetAllPostsHandler(queries).ServeHTTP)
	apiRouter.Post("/posts", middlewares.MiddlewareAuth(queries, nil,
		func(w http.ResponseWriter, r *http.Request, contributor database.Contributor) {
			handlers.CreatePostHandler(queries, contributor).ServeHTTP(w, r)
		}, nil, "contributor"))
	apiRouter.Put("/posts/{id}", middlewares.MiddlewareAuth(queries, nil,
		func(w http.ResponseWriter, r *http.Request, contributor database.Contributor) {
			handlers.UpdatePostHandler(queries, contributor).ServeHTTP(w, r)
		}, nil, "contributor"))
	apiRouter.Delete("/posts/{id}", middlewares.MiddlewareAuth(queries, nil,
		func(w http.ResponseWriter, r *http.Request, contributor database.Contributor) {
			handlers.DeletePostHandler(queries, contributor).ServeHTTP(w, r)
		}, nil, "contributor"))
	apiRouter.Get("/posts/{slug}", middlewares.MiddlewareModeratorOrUser(queries,
		func(w http.ResponseWriter, r *http.Request, u database.User) {
			handlers.GetPostBySlugHandler(queries, u, database.Moderator{}).ServeHTTP(w, r)
		},
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.GetPostBySlugHandler(queries, database.User{}, m).ServeHTTP(w, r)
		}))

	// Post Interactions Routes
	apiRouter.Post("/posts/{postID}/upvotes", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.InsertUpvoteHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))
	apiRouter.Delete("/posts/{postID}/upvotes", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.DeleteUpvoteHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))
	apiRouter.Get("/upvotes", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.GetUpvotesHandlerByUser(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	// Comments Routes
	apiRouter.Post("/posts/{postID}/comments", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.CreateCommentHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))
	apiRouter.Delete("/posts/{postID}/comments/{commentID}", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.DeleteCommentHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))
	apiRouter.Patch("/posts/{postID}/comments/{commentID}", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.UpdateCommentHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))
	apiRouter.Get("/posts/{postSlug}/comments", middlewares.MiddlewareModeratorOrUser(queries,
		func(w http.ResponseWriter, r *http.Request, u database.User) {
			handlers.GetAllCommentsByPostHandler(queries).ServeHTTP(w, r)
		},
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.GetAllCommentsByPostHandler(queries).ServeHTTP(w, r)
		}))

	// Follow Routes
	apiRouter.Post("/follow", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.CreateFollowHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))
	apiRouter.Delete("/follow/{id}", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.DeleteFollowHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))
	apiRouter.Get("/users/{username}/following", middlewares.MiddlewareModeratorOrUser(queries,
		func(w http.ResponseWriter, r *http.Request, u database.User) {
			handlers.GetFollowingListByIDHandler(queries).ServeHTTP(w, r)
		},
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.GetFollowingListByIDHandler(queries).ServeHTTP(w, r)
		}))

	// Feed Route
	apiRouter.Get("/feed", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlers.GetFeedHandler(queries, user).ServeHTTP(w, r)
		}, nil, nil, "user"))

	// Saved Posts Routes
	apiRouter.Post("/saved-posts", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, u database.User) {
			handlers.CreateSavePost(queries, u).ServeHTTP(w, r)
		}, nil, nil, "user"))
	apiRouter.Delete("/saved-posts/{id}", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, u database.User) {
			handlers.DeleteSavedPost(queries, u).ServeHTTP(w, r)
		}, nil, nil, "user"))
	apiRouter.Get("/saved-posts/{username}", middlewares.MiddlewareModeratorOrUser(queries,
		func(w http.ResponseWriter, r *http.Request, u database.User) {
			handlers.GetSavedPosts(queries).ServeHTTP(w, r)
		},
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.GetSavedPosts(queries).ServeHTTP(w, r)
		}))

	// Profile Routes
	apiRouter.Get("/profile/{username}", middlewares.MiddlewareModeratorOrUser(queries,
		func(w http.ResponseWriter, r *http.Request, u database.User) {
			handlers.GetProfileDataHandler(queries, u, database.Moderator{}).ServeHTTP(w, r)
		},
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.GetProfileDataHandler(queries, database.User{}, m).ServeHTTP(w, r)
		}))
	apiRouter.Get("/profile/{username}/posts", middlewares.MiddlewareModeratorOrUser(queries,
		func(w http.ResponseWriter, r *http.Request, u database.User) {
			handlers.GetContributorProfilePostsHandler(queries).ServeHTTP(w, r)
		},
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.GetContributorProfilePostsHandler(queries).ServeHTTP(w, r)
		}))
	apiRouter.Put("/profile/update", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, u database.User) {
			handlers.UpdateUserHandler(queries, u).ServeHTTP(w, r)
		}, nil, nil, "user"))
	apiRouter.Get("/profile/reports", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, u database.User) {
			handlers.GetResolvedReportsWithSuspensionHandler(queries, u).ServeHTTP(w, r)
		}, nil, nil, "user"))

	// Admin Routes
	apiRouter.Post("/admin/login", handlers.LoginModeratorController(queries).ServeHTTP)
	apiRouter.Post("/admin/create", middlewares.MiddlewareAuth(queries, nil, nil,
		func(w http.ResponseWriter, r *http.Request, moderator database.Moderator) {
			handlers.CreateModeratorHandler(queries, moderator).ServeHTTP(w, r)
		}, "moderator"))
	apiRouter.Get("/admin/moderators", middlewares.MiddlewareAuth(queries, nil, nil,
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.GetAllModerators(queries).ServeHTTP(w, r)
		}, "moderator"))

	// Contributor Application Routes
	apiRouter.Post("/contributor-applications", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, u database.User) {
			handlers.CreateContributorApplication(queries, u).ServeHTTP(w, r)
		}, nil, nil, "user"))
	apiRouter.Get("/admin/contributor-applications", middlewares.MiddlewareAuth(queries, nil, nil,
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.GetContributorApplications(queries, m).ServeHTTP(w, r)
		}, "moderator"))
	apiRouter.Get("/admin/contributor-applications/{id}", middlewares.MiddlewareAuth(queries, nil, nil,
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.GetContributorApplicationByID(queries, m).ServeHTTP(w, r)
		}, "moderator"))
	apiRouter.Put("/admin/contributor-applications/{id}/status", middlewares.MiddlewareAuth(queries, nil, nil,
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.UpdateContributorApplication(queries, m).ServeHTTP(w, r)
		}, "moderator"))

	// Reports Routes
	apiRouter.Post("/reports", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, u database.User) {
			handlers.CreateReportHandler(queries, u).ServeHTTP(w, r)
		}, nil, nil, "user"))
	apiRouter.Get("/admin/contributors/reports", middlewares.MiddlewareAuth(queries, nil, nil,
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.GetReportedContributorsHandler(queries, m).ServeHTTP(w, r)
		}, "moderator"))
	apiRouter.Get("/admin/users/reports", middlewares.MiddlewareAuth(queries, nil, nil,
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.GetReportedUserHandler(queries, m).ServeHTTP(w, r)
		}, "moderator"))
	apiRouter.Put("/admin/reports/{reportID}/status", middlewares.MiddlewareAuth(queries, nil, nil,
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.UpdateReportStatusHandler(queries, m).ServeHTTP(w, r)
		}, "moderator"))

	// Appeals Routes
	apiRouter.Post("/appeals", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, u database.User) {
			handlers.CreateAppealHandler(queries, u).ServeHTTP(w, r)
		}, nil, nil, "user"))
	apiRouter.Get("/admin/appeals", middlewares.MiddlewareAuth(queries, nil, nil,
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.GetAppealsHandler(queries, m).ServeHTTP(w, r)
		}, "moderator"))
	apiRouter.Get("/admin/contributors/appeals", middlewares.MiddlewareAuth(queries, nil, nil,
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.GetContributorsAppeals(queries).ServeHTTP(w, r)
		}, "moderator"))
	apiRouter.Get("/admin/users/appeals", middlewares.MiddlewareAuth(queries, nil, nil,
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.GetUsersAppeals(queries).ServeHTTP(w, r)
		}, "moderator"))
	apiRouter.Get("/admin/appeals/{appealID}", middlewares.MiddlewareAuth(queries, nil, nil,
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.GetAppealByIDHandler(queries, m).ServeHTTP(w, r)
		}, "moderator"))
	apiRouter.Put("/admin/appeals/{appealID}/status", middlewares.MiddlewareAuth(queries, nil, nil,
		func(w http.ResponseWriter, r *http.Request, m database.Moderator) {
			handlers.UpdateAppealStatus(queries, m).ServeHTTP(w, r)
		}, "moderator"))

	// Search Routes
	apiRouter.Get("/search/posts", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, u database.User) {
			handlers.SearchPostsHandler(queries).ServeHTTP(w, r)
		}, nil, nil, "user"))
	apiRouter.Get("/search/users", middlewares.MiddlewareAuth(queries,
		func(w http.ResponseWriter, r *http.Request, u database.User) {
			handlers.SearchUsersHandler(queries).ServeHTTP(w, r)
		}, nil, nil, "user"))

	return r
}
