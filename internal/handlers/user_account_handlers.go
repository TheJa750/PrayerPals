package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/TheJa750/PrayerPals/internal/auth"
	"github.com/TheJa750/PrayerPals/internal/database"
)

func (a *APIConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body for user data
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

	// Add user to the database
	userData, err := a.DBQueries.CreateUser(r.Context(), database.CreateUserParams{
		Username:       userReq.Username,
		Email:          userReq.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Create JSON response
	jsonUser := User{
		ID:       userData.ID,
		Username: userData.Username,
		Email:    userData.Email,
	}

	err = CreateJSONResponse(jsonUser, w, http.StatusCreated)
	if err != nil {
		log.Printf("Error creating JSON response: %v", err)
		return
	}

	log.Printf("User %s created successfully", userData.Username)
}

func (a *APIConfig) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
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

func (a *APIConfig) RefreshJWTHandler(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Invalid token format", http.StatusBadRequest)
		return
	}

	// Check user refresh token
	refreshToken, err := a.DBQueries.GetUserByToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Assign user ID from the refresh token
	userID := refreshToken.UserID

	// Issue new tokens
	accessToken, err := auth.MakeJWT(userID, a.JWTSecret, 1800*time.Second)
	if err != nil {
		log.Printf("Error generating access token for user %v: %v", userID, err)
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}

	jsonResponse := UserLoggedIn{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: token,
	}

	err = CreateJSONResponse(jsonResponse, w, http.StatusOK)
	if err != nil {
		log.Printf("Error creating JSON response: %v", err)
		return
	}

	log.Printf("Tokens refreshed for user %v", userID)
}
