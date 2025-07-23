package handlers

import (
	"github.com/TheJa750/PrayerPals/internal/database"
)

type JsonError struct {
	Message string `json:"error"`
}

type APIConfig struct {
	DBQueries *database.Queries
	JWTSecret string
}

type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
