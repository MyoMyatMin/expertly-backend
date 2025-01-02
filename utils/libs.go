package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/MyoMyatMin/expertly-backend/pkg/database"
)

func GenerateUniqueUsername(name string, db *database.Queries, r *http.Request) (string, error) {
	username := strings.ReplaceAll(strings.ToLower(name), " ", "")

	exists, err := checkUsernameExists(username, db, r)
	if err != nil {
		return "", fmt.Errorf("failed to check username availability: %v", err)
	}

	if exists {
		suffixNumber := 1
		uniqueUsername := fmt.Sprintf("%s%d", username, suffixNumber)
		for {
			exists, err = checkUsernameExists(uniqueUsername, db, r)
			if err != nil {
				return "", fmt.Errorf("failed to check username availability: %v", err)
			}
			if !exists {
				return uniqueUsername, nil
			}
			suffixNumber++
			uniqueUsername = fmt.Sprintf("%s%d", username, suffixNumber)
		}
	}

	return username, nil
}

func checkUsernameExists(username string, db *database.Queries, r *http.Request) (bool, error) {
	_, err := db.GetUserByUsername(r.Context(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GenerateUniqueSlug(title string, db *database.Queries, r *http.Request) (string, error) {
	slug := strings.ReplaceAll(strings.ToLower(title), " ", "-")

	exists, err := checkSlugExists(slug, db, r)
	if err != nil {
		return "", fmt.Errorf("failed to check slug availability: %v", err)
	}

	if exists {
		suffixNumber := 1
		uniqueSlug := fmt.Sprintf("%s%d", slug, suffixNumber)
		for {
			exists, err = checkSlugExists(uniqueSlug, db, r)
			if err != nil {
				return "", fmt.Errorf("failed to check slug availability: %v", err)
			}
			if !exists {
				return uniqueSlug, nil
			}
			suffixNumber++
			uniqueSlug = fmt.Sprintf("%s%d", slug, suffixNumber)
		}
	}

	return slug, nil
}

func checkSlugExists(slug string, db *database.Queries, r *http.Request) (bool, error) {
	_, err := db.GetPostBySlug(r.Context(), slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
