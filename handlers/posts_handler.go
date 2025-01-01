package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func GetAllPostsHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		posts, err := db.ListPosts(r.Context())
		if err != nil {
			http.Error(w, "Couldn't get posts", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // 200
		json.NewEncoder(w).Encode(posts)
	})
}

func CreatePostHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Title   string   `json:"title"`
			Content string   `json:"content"`
			Images  []string `json:"images"`
		}

		var params parameters

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest) // 400
			return
		}

		imagesToDelete := findImagesToDelete(params.Content, params.Images)
		for _, imageURL := range imagesToDelete {
			publicID := extractPublicID(imageURL)
			if publicID != "" {
				err := deleteImagesFromCloudinary(publicID)
				if err != nil {
					fmt.Printf("Failed to delete image: %s, error: %v\n", imageURL, err)
				}
			}
		}

		post, err := db.CreatePost(r.Context(), database.CreatePostParams{
			ID:      uuid.New(),
			Title:   params.Title,
			Content: params.Content,
			UserID:  user.ID,
		})
		if err != nil {
			http.Error(w, "Couldn't create post", http.StatusInternalServerError) // 500
			return
		}

		w.WriteHeader(http.StatusCreated) // 201
		json.NewEncoder(w).Encode(post)
	})
}

func GetPostByIDHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "id")

		postUUID, err := uuid.Parse(postID)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest) // 400
			return
		}

		post, err := db.GetPost(r.Context(), postUUID)
		if err != nil {
			http.Error(w, "Couldn't get post", http.StatusNotFound) // 404
			return
		}

		w.WriteHeader(http.StatusOK) // 200
		json.NewEncoder(w).Encode(post)
	})
}

func UpdatePostHandler(db *database.Queries, user database.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "id")

		postUUID, err := uuid.Parse(postID)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest) // 400
			return
		}

		type parameters struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}

		var params parameters

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest) // 400
			return
		}

		post, err := db.UpdatePost(r.Context(), database.UpdatePostParams{
			ID:      postUUID,
			Title:   params.Title,
			Content: params.Content,
		})
		if err != nil {
			http.Error(w, "Couldn't update post", http.StatusInternalServerError) // 500
			return
		}

		w.WriteHeader(http.StatusOK) // 200
		json.NewEncoder(w).Encode(post)
	})
}

func findImagesToDelete(content string, images []string) []string {
	var imagesToDelete []string
	for _, image := range images {
		if !strings.Contains(content, image) {
			imagesToDelete = append(imagesToDelete, image)
		}
	}
	return imagesToDelete
}

func extractPublicID(imgURL string) string {
	parts := strings.Split(imgURL, "/")
	if len(parts) < 2 {
		return ""
	}

	publicIDWithExt := parts[len(parts)-1]
	publicID := strings.Split(publicIDWithExt, ".")[0]
	return publicID
}

func deleteImagesFromCloudinary(publicID string) error {
	godotenv.Load(".env")
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return fmt.Errorf("failed to initialize Cloudinary: %v", err)
	}
	_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: publicID})

	return err
}
