package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/MyoMyatMin/expertly-backend/routes"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	db     *sql.DB
	dbOnce sync.Once
)

func connectToDB() *sql.DB {
	dbOnce.Do(func() {
		godotenv.Load(".env")
		dbURL := os.Getenv("DB_URL")
		if dbURL == "" {
			log.Fatal("DB_URL is not set")
		}

		var err error
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("Error opening database: %q", err)
		}

		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)
		db.SetConnMaxIdleTime(5 * time.Minute)
		if err := db.Ping(); err != nil {
			log.Fatalf("Error connecting to DB: %v", err)
		}
	})
	return db
}

func Handler(w http.ResponseWriter, r *http.Request) {
	db := connectToDB()
	router := routes.SetUpRoutes(db)
	router.ServeHTTP(w, r)
}
