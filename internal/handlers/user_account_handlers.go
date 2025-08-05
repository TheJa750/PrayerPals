package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/TheJa750/PrayerPals/internal/auth"
	"github.com/TheJa750/PrayerPals/internal/validation"
)

func (a *APIConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body for user data
	userReq, err := ParseJSON[UserRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate user input
	if userReq.Username == "" || userReq.Email == "" || userReq.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	validEmail := validation.ValidateEmail(userReq.Email)
	if !validEmail.IsValid {
		http.Error(w, strings.Join(validEmail.Errors, ", "), http.StatusBadRequest)
		return
	}

	user, err := a.createUser(r.Context(), userReq)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	err = CreateJSONResponse(user, w, http.StatusCreated)
	if err != nil {
		log.Printf("Error creating JSON response: %v", err)
		return
	}

	log.Printf("User %s created successfully", user.Username)
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
	userData, err := a.DBQueries.GetUserIDByEmail(r.Context(), strings.ToLower(loginReq.Email))
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
		ID:       userData.ID,
		Username: userData.Username,
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Use true in production (HTTPS)
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Use true in production (HTTPS)
		SameSite: http.SameSiteStrictMode,
	})

	err = CreateJSONResponse(jsonUser, w, http.StatusOK)
	if err != nil {
		log.Printf("Error creating JSON response: %v", err)
		return
	}

	log.Printf("User %s logged in successfully", userData.Username)
}

func (a *APIConfig) RefreshJWTHandler(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from header
	token, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Missing refresh token", http.StatusUnauthorized)
		return
	}

	// Check user refresh token
	refreshToken, err := a.DBQueries.GetUserByToken(r.Context(), token.Value)
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
		ID: userID,
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Use true in production (HTTPS)
		SameSite: http.SameSiteStrictMode,
	})

	err = CreateJSONResponse(jsonResponse, w, http.StatusOK)
	if err != nil {
		log.Printf("Error creating JSON response: %v", err)
		return
	}

	log.Printf("Tokens refreshed for user %v", userID)
}

func (a *APIConfig) LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	// Clear cookies to log out user
	clearCookie := func(name string) {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // Set to true in production
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1, // Delete immediately
		})
	}
	clearCookie("access_token")
	clearCookie("refresh_token")

	// Revoke refresh token in the database
	token, err := r.Cookie("refresh_token")
	if err == nil {
		_ = a.DBQueries.RevokeUserToken(r.Context(), token.Value)
	}

	log.Println("User logged out successfully")
	w.WriteHeader(http.StatusNoContent)
}
