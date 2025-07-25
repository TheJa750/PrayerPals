package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

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

func (a *APIConfig) PromoteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body for user promotion details
	promoteReq, err := ParseJSON[PromoteUserRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if the user is admin in group
	userRole, err := a.DBQueries.GetUserGroupRole(r.Context(), database.GetUserGroupRoleParams{
		UserID:  userID,
		GroupID: promoteReq.GroupID,
	})
	if err != nil {
		http.Error(w, "Failed to get user role in group", http.StatusInternalServerError)
		return
	}

	if userRole != "admin" {
		http.Error(w, "Only admins can promote users", http.StatusForbidden)
		return
	}

	// Perform checks before promoting user
	role := strings.ToLower(promoteReq.Role)
	err = a.promoteUserChecks(r.Context(), promoteReq.TargetUserID, promoteReq.GroupID, role)
	if err != nil {
		if errors.Is(err, ErrUserNotMember) || errors.Is(err, ErrInvalidRole) || errors.Is(err, ErrUserHasRole) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to promote user", http.StatusInternalServerError)
		return
	}

	// Update user role in the database
	err = a.DBQueries.AdjustUserGroupRole(r.Context(), database.AdjustUserGroupRoleParams{
		UserID:  promoteReq.TargetUserID,
		GroupID: promoteReq.GroupID,
		Role:    role,
	})
	if err != nil {
		http.Error(w, "Failed to promote user", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusNoContent)
	log.Printf("User %v promoted to %s in group %v", promoteReq.TargetUserID, role, promoteReq.GroupID)

}
