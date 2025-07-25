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

	// Server Admin Handlers
	mux.HandleFunc("POST /admin/reset", cfg.ResetDatabase)
	mux.HandleFunc("POST /admin/reset/users", cfg.ResetUsersOnly)
	mux.HandleFunc("POST /admin/reset/groups", cfg.ResetGroupsOnly)

	// Generic API Handlers
	mux.HandleFunc("GET /api/health", handlers.HealthCheck)

	// User Account Handlers
	mux.HandleFunc("POST /api/users", cfg.CreateUserHandler)
	mux.HandleFunc("POST /api/login", cfg.LoginUserHandler)

	// User Functions Handlers
	mux.HandleFunc("POST /api/groups/join", cfg.JoinGroupHandler)    // Expecting query parameter ?group_id=UUID
	mux.HandleFunc("DELETE /api/groups/join", cfg.LeaveGroupHandler) // Expecting query parameter ?group_id=UUID
	mux.HandleFunc("GET /api/groups", cfg.GetGroupsForFeed)          // Fetch groups for user feed

	// Group Handlers
	mux.HandleFunc("POST /api/groups", cfg.CreateGroupHandler)
	mux.HandleFunc("PUT /api/groups/promote", cfg.PromoteUserHandler)

	// Post Handlers
	mux.HandleFunc("POST /api/posts", cfg.CreatePostHandler)
	mux.HandleFunc("DELETE /api/posts", cfg.DeletePostHandler) // Expecting JSON
	mux.HandleFunc("POST /api/comments", cfg.CreateCommentHandler)

	log.Fatal(svr.ListenAndServe())
}
