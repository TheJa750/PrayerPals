package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/TheJa750/PrayerPals/internal/auth"
	"github.com/TheJa750/PrayerPals/internal/database"
)

func (a *APIConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic for creating a user
	userReq, err := ParseJSON[UserRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if userReq.Username == "" || userReq.Email == "" || userReq.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(userReq.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = a.DBQueries.CreateUser(r.Context(), database.CreateUserParams{
		Username:       userReq.Username,
		Email:          userReq.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *APIConfig) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Handler logic for user login
	loginReq, err := ParseJSON[UserRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if loginReq.Email == "" || loginReq.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Retrieve user by email
	userData, err := a.DBQueries.GetUserIDByEmail(r.Context(), loginReq.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check password
	if !auth.CheckPasswordHash(loginReq.Password, userData.HashedPassword) {
		log.Printf("Password mismatch for user %s, error: %v", userData.Username, err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Issue tokens
	accessToken, refreshToken, err := a.issueTokens(userData, a.JWTSecret, 1800*time.Second, r.Context())
	if err != nil {
		log.Printf("Error issuing tokens for user %s: %v", userData.Username, err)
		http.Error(w, "Error issuing tokens", http.StatusInternalServerError)
		return
	}

	jsonUser := UserLoggedIn{
		ID:           userData.ID,
		Username:     userData.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	err = CreateJSONResponse(jsonUser, w, http.StatusOK)
	if err != nil {
		log.Printf("Error creating JSON response: %v", err)
		return
	}

	log.Printf("User %s logged in successfully", userData.Username)
}
