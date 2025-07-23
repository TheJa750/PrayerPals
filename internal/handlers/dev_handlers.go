package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func (a *APIConfig) ResetDatabase(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		http.Error(w, "This endpoint is only available in development mode", http.StatusForbidden)
		return
	}

	err := a.DBQueries.ResetUsers(r.Context())
	if err != nil {
		http.Error(w, "Error resetting users", http.StatusInternalServerError)
		return
	}

	err = a.DBQueries.ResetGroups(r.Context())
	if err != nil {
		http.Error(w, "Error resetting groups", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Users and Groups reset successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (a *APIConfig) ResetUsersOnly(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		http.Error(w, "This endpoint is only available in development mode", http.StatusForbidden)
		return
	}

	err := a.DBQueries.ResetUsers(r.Context())
	if err != nil {
		http.Error(w, "Error resetting users", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Users reset successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (a *APIConfig) ResetGroupsOnly(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		http.Error(w, "This endpoint is only available in development mode", http.StatusForbidden)
		return
	}

	err := a.DBQueries.ResetGroups(r.Context())
	if err != nil {
		http.Error(w, "Error resetting groups", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Groups reset successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
