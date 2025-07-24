package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TheJa750/PrayerPals/internal/database"
	"github.com/TheJa750/PrayerPals/internal/handlers"
	"github.com/TheJa750/PrayerPals/internal/middleware"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	mux := http.NewServeMux()

	svr := http.Server{
		Addr:              ":8080",
		Handler:           middleware.LoggingMiddleware(mux),
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Println("Starting server on :8080")

	cfg := handlers.APIConfig{
		DBQueries: database.New(db),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	fileHandler := http.StripPrefix("/app/", http.FileServer(http.Dir("./internal/assets/")))

	mux.Handle("/app/", fileHandler)

	//Admin Handlers
	mux.HandleFunc("POST /admin/reset", cfg.ResetDatabase)
	mux.HandleFunc("POST /admin/reset/users", cfg.ResetUsersOnly)
	mux.HandleFunc("POST /admin/reset/groups", cfg.ResetGroupsOnly)

	// API Handlers
	mux.HandleFunc("GET /api/health", handlers.HealthCheck)
	mux.HandleFunc("POST /api/users", cfg.CreateUserHandler)
	mux.HandleFunc("POST /api/login", cfg.LoginUser)
	mux.HandleFunc("POST /api/groups", cfg.CreateGroup)

	log.Fatal(svr.ListenAndServe())
}
