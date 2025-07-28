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
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()

	svr := http.Server{
		Addr:              ":8080",
		Handler:           middleware.LoggingMiddleware(router),
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Println("Starting server on :8080")

	cfg := handlers.APIConfig{
		DBQueries: database.New(db),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	// File handler
	router.PathPrefix("/app/").Handler(http.StripPrefix("/app/", http.FileServer(http.Dir("./internal/assets/"))))

	// Server Admin Handlers
	router.HandleFunc("/admin/reset", cfg.ResetDatabase).Methods("POST")
	router.HandleFunc("/admin/reset/users", cfg.ResetUsersOnly).Methods("POST")
	router.HandleFunc("/admin/reset/groups", cfg.ResetGroupsOnly).Methods("POST")

	// Generic API Handlers
	router.HandleFunc("/api/health", handlers.HealthCheck).Methods("GET")

	// User Account Handlers
	router.HandleFunc("/api/users", cfg.CreateUserHandler).Methods("POST")
	router.HandleFunc("/api/login", cfg.LoginUserHandler).Methods("POST")
	router.HandleFunc("/api/refresh", cfg.RefreshJWTHandler).Methods("POST")

	// User Functions Handlers
	router.HandleFunc("/api/groups/join/{group_id}", cfg.JoinGroupHandler).Methods("POST")
	router.HandleFunc("/api/groups/join/{group_id}", cfg.LeaveGroupHandler).Methods("DELETE")
	router.HandleFunc("/api/groups", cfg.GetGroupsForFeed).Methods("GET")

	// Group Handlers
	router.HandleFunc("/api/groups", cfg.CreateGroupHandler).Methods("POST")
	router.HandleFunc("/api/groups/promote", cfg.PromoteUserHandler).Methods("PUT")
	router.HandleFunc("/api/groups/{group_id}/posts", cfg.GetPostFeedHandler).Methods("GET") // Expecting query parameters ?limit=10&offset=0
	router.HandleFunc("/api/groups/{group_id}", cfg.DeleteGroupHandler).Methods("DELETE")

	// Post Handlers
	router.HandleFunc("/api/posts", cfg.CreatePostHandler).Methods("POST")
	router.HandleFunc("/api/posts", cfg.DeletePostHandler).Methods("DELETE") // Expecting JSON
	router.HandleFunc("/api/comments", cfg.CreateCommentHandler).Methods("POST")
	router.HandleFunc("/api/comments/{post_id}", cfg.GetCommentsForPostHandler).Methods("GET")

	log.Fatal(svr.ListenAndServe())
}
