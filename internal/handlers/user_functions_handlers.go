package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/TheJa750/PrayerPals/internal/database"
	"github.com/google/uuid"
)

func (a *APIConfig) JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request URL for group ID
	groupID, err := uuid.Parse(r.URL.Query().Get("group_id"))
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Add user to the group in database
	err = a.DBQueries.AddUserToGroup(r.Context(), database.AddUserToGroupParams{
		UserID:  userID,
		GroupID: groupID,
		Role:    "member",
	})
	if err != nil {
		http.Error(w, "Failed to join group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *APIConfig) LeaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request URL for group ID
	groupID, err := uuid.Parse(r.URL.Query().Get("group_id"))
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Verify user can leave the group
	// This will check if the user is a member, if they are the only admin,
	// and if they are the last member of the group
	// If any check fails, it will return an error
	canLeaveErr := a.leaveGroupChecks(r.Context(), userID, groupID)
	if canLeaveErr != nil {
		if errors.Is(canLeaveErr, ErrUserNotMember) ||
			errors.Is(canLeaveErr, ErrUserIsOnlyAdmin) ||
			errors.Is(canLeaveErr, ErrUserIsLastMember) {
			http.Error(w, canLeaveErr.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to leave group", http.StatusInternalServerError)
		return
	}

	// Remove user from the group in database
	err = a.DBQueries.RemoveUserFromGroup(r.Context(), database.RemoveUserFromGroupParams{
		UserID:  userID,
		GroupID: groupID,
	})
	if err != nil {
		http.Error(w, "Failed to leave group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("User %v left group %v", userID, groupID)
}
