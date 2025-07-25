package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/TheJa750/PrayerPals/internal/database"
	"github.com/google/uuid"
)

func (a *APIConfig) CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic for creating a group
	// Parse JSON request body
	groupReq, err := ParseJSON[GroupRequest](r)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate name field
	if groupReq.Name == "" {
		log.Println("Missing group name")
		http.Error(w, "Missing required field", http.StatusBadRequest)
		return
	}

	// Setting up query parameters
	// Description can be null, so we use sql.NullString
	valid := true
	if groupReq.Description == "" {
		valid = false
	}

	description := sql.NullString{
		String: groupReq.Description,
		Valid:  valid,
	}
	// Validate JWT and get user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Error getting user ID from token: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create the group in the database
	group, err := a.DBQueries.CreateGroup(r.Context(), database.CreateGroupParams{
		Name:        groupReq.Name,
		Description: description,
		OwnerID: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	})
	if err != nil {
		http.Error(w, "Error creating group", http.StatusInternalServerError)
		return
	}

	// Add user to group
	err = a.DBQueries.AddUserToGroup(r.Context(), database.AddUserToGroupParams{
		UserID:  userID,
		GroupID: group.ID,
		Role:    "admin", // Group creator is admin
	})
	if err != nil {
		log.Printf("Error adding user to group: %v", err)
		http.Error(w, "Error adding user to group", http.StatusInternalServerError)
		return
	}

	// Create JSON response
	jsonGroup := Group{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description.String,
		OwnerID:     group.OwnerID.UUID,
	}

	if err := CreateJSONResponse(jsonGroup, w, http.StatusCreated); err != nil {
		log.Printf("Error creating JSON response: %v", err)
		return
	}

	log.Printf("Group created successfully: %s", group.Name)
}
