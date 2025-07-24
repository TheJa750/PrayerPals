package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

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

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type UserLoggedIn struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

type GroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Group struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     uuid.UUID `json:"owner_id"`
}

func ParseJSON[T any](r *http.Request) (T, error) {
	var data T

	if r.Body == nil {
		return data, errors.New("request body is empty")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		return data, errors.New("invalid request body")
	}

	return data, nil
}

func CreateJSONResponse[T any](data T, w http.ResponseWriter, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return err
	}
	return nil
}
