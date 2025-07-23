package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

func main() {
	/*
		godotenv.Load()
		dbURL := os.Getenv("DB_URL")

		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("Error connecting to the database: %v", err)
		}
	*/

	mux := http.NewServeMux()

	svr := http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Println("Starting server on :8080")

	mux.Handle("/", http.FileServer(http.Dir(".")))

	log.Fatal(svr.ListenAndServe())
}
