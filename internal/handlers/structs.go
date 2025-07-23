package handlers

import (
	"github.com/TheJa750/PrayerPals/internal/database"
	"github.com/google/uuid"
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

type UserLoggedIn struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}
