package main

import (
	"log"
	"net/http"
	"time"

	"github.com/TheJa750/PrayerPals/internal/handlers"
	"github.com/TheJa750/PrayerPals/internal/middleware"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

func main() {
	mux := http.NewServeMux()

	svr := http.Server{
		Addr:              ":8080",
		Handler:           middleware.LoggingMiddleware(mux),
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Println("Starting server on :8080")

	fileHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", fileHandler)
	mux.HandleFunc("/health", handlers.HealthCheck)

	log.Fatal(svr.ListenAndServe())
}

/*
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
*/
