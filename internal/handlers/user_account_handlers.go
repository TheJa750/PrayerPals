package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/TheJa750/PrayerPals/internal/auth"
	"github.com/TheJa750/PrayerPals/internal/database"
)

func (a *APIConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic for creating a user
	var userReq UserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

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
	var loginReq UserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if loginReq.Email == "" || loginReq.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Retrieve user by email
	id, err := a.DBQueries.GetUserIDByEmail(r.Context(), loginReq.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	//Get user information from id
	userData, err := a.DBQueries.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Error retrieving user data", http.StatusInternalServerError)
		return
	}

	// Check password
	if !auth.CheckPasswordHash(loginReq.Password, userData.HashedPassword) {
		log.Printf("Password mismatch for user %s, error: %v", userData.Username, err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	access_token, err := auth.MakeJWT(userData.ID, a.JWTSecret, 1800*time.Second)
	if err != nil {
		log.Printf("Error creating JWT: %v", err)
		http.Error(w, "Error creating access token", http.StatusInternalServerError)
		return
	}

	rtString, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error creating refresh token: %v", err)
		http.Error(w, "Error creating refresh token", http.StatusInternalServerError)
		return
	}

	_, err = a.DBQueries.CreateUserToken(r.Context(), database.CreateUserTokenParams{
		UserID: userData.ID,
		Token:  rtString,
	})
	if err != nil {
		log.Printf("Error saving refresh token: %v", err)
		http.Error(w, "Error saving refresh token", http.StatusInternalServerError)
		return
	}

	jsonUser := UserLoggedIn{
		ID:           userData.ID,
		Username:     userData.Username,
		AccessToken:  access_token,
		RefreshToken: rtString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(jsonUser); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	log.Printf("User %s logged in successfully", userData.Username)
}
