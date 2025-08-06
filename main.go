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
		Handler:           middleware.CorsMiddleware(middleware.LoggingMiddleware(router)),
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
	router.HandleFunc("/api/users", cfg.CreateUserHandler).Methods("POST")       // Expecting JSON body for username/email/password
	router.HandleFunc("/api/users/update", cfg.UpdateUserHandler).Methods("PUT") // Expecting JSON body for username/password (1 only)
	router.HandleFunc("/api/login", cfg.LoginUserHandler).Methods("POST")        // Expecting JSON body for email/password
	router.HandleFunc("/api/refresh", cfg.RefreshJWTHandler).Methods("POST")
	router.HandleFunc("/api/logout", cfg.LogoutUserHandler).Methods("POST")

	// User Functions Handlers
	router.HandleFunc("/api/groups/{invite_code}/join", cfg.JoinGroupHandler).Methods("POST")
	router.HandleFunc("/api/groups/{group_id}/join", cfg.LeaveGroupHandler).Methods("DELETE")
	router.HandleFunc("/api/groups", cfg.GetGroupsForFeed).Methods("GET")

	// Group Handlers
	router.HandleFunc("/api/groups", cfg.CreateGroupHandler).Methods("POST")                                     // Expecting JSON body for name/description
	router.HandleFunc("/api/groups/{group_id}", cfg.GetGroupInfoHandler).Methods("GET")                          // Expecting group_id in URL
	router.HandleFunc("/api/groups/{group_id}/members/{user_id}/promote", cfg.PromoteUserHandler).Methods("PUT") // Expecting JSON body for new role
	router.HandleFunc("/api/groups/{group_id}/posts", cfg.GetPostFeedHandler).Methods("GET")                     // Expecting query parameters ?limit=10&offset=0
	router.HandleFunc("/api/groups/{group_id}", cfg.DeleteGroupHandler).Methods("DELETE")
	router.HandleFunc("/api/groups/{group_id}/members/{user_id}/moderate", cfg.ModerateUserHandler).Methods("PUT") // Expecting JSON body for action and reason
	router.HandleFunc("/api/groups/{invite_code}", cfg.GroupFromInviteCodeHandler).Methods("GET")
	router.HandleFunc("/api/groups/{group_id}/members", cfg.GetGroupMembersHandler).Methods("GET")            // Expecting group_id in URL
	router.HandleFunc("/api/groups/{group_id}/members/{user_id}", cfg.GetUserGroupRoleHandler).Methods("GET") // Expecting group_id and user_id in URL

	// Post Handlers
	router.HandleFunc("/api/groups/{group_id}/posts", cfg.CreatePostHandler).Methods("POST") // Expecting JSON body for post content
	router.HandleFunc("/api/groups/{group_id}/posts/{post_id}", cfg.DeletePostHandler).Methods("DELETE")
	router.HandleFunc("/api/groups/{group_id}/posts/{post_id}/comments", cfg.CreateCommentHandler).Methods("POST") // Expecting JSON body for comment content
	router.HandleFunc("/api/groups/{group_id}/posts/{post_id}/comments", cfg.GetCommentsForPostHandler).Methods("GET")

	log.Fatal(svr.ListenAndServe())
}
