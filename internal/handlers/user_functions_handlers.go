package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/TheJa750/PrayerPals/internal/database"
)

func (a *APIConfig) JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request URL for invite code
	invCode, err := parseInviteCodePathParam(r, "invite_code")
	if err != nil {
		http.Error(w, "Invalid group invite code", http.StatusBadRequest)
		return
	}

	// Fetch group by invite code
	group, err := a.DBQueries.GetGroupByInviteCode(r.Context(), invCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Group not found for invite code", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching group by invite code: %v", err)
		http.Error(w, "Failed to fetch group by invite code", http.StatusInternalServerError)
		return
	}

	// Verify user is not already a member
	isMember, err := a.verifyUserInGroup(r.Context(), userID, group.ID)
	if err != nil {
		log.Printf("Error verifying user in group: %v", err)
		http.Error(w, "Failed to verify group membership", http.StatusInternalServerError)
		return
	}
	if isMember {
		http.Error(w, "User is already a member of the group", http.StatusConflict)
		return
	}

	// Add user to the group in database
	err = a.joinGroup(r.Context(), userID, group.ID, "member")
	if err != nil {
		log.Printf("Error adding user to group: %v", err)
		http.Error(w, "Error adding user to group", http.StatusInternalServerError)
		return
	}

	// Success response with group info
	response := map[string]interface{}{
		"message":    "Successfully joined group",
		"group_name": group.Name,
		"group_id":   group.ID,
	}
	if err := CreateJSONResponse(response, w, http.StatusOK); err != nil {
		log.Printf("Error creating JSON response: %v", err)
		return
	}

	log.Printf("User %v joined group %v successfully", userID, group.ID)
}

func (a *APIConfig) LeaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request URL for group ID
	groupID, err := parseUUIDPathParam(r, "group_id")
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

func (a *APIConfig) GetGroupsForFeed(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch groups for the user from the database
	groupIDs, err := a.DBQueries.GetGroupsForUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch groups", http.StatusInternalServerError)
		return
	}

	// Create slice of JSON groups for response
	jsonGroups := make([]Group, len(groupIDs))
	for i, groupID := range groupIDs {
		group, err := a.DBQueries.GetGroupByID(r.Context(), groupID)
		if err != nil {
			log.Printf("Failed to fetch group %v: %v", groupID, err)
		}

		jsonGroups[i] = Group{
			ID:          group.ID,
			Name:        group.Name,
			OwnerID:     group.OwnerID.UUID,
			Description: group.Description.String,
		}
	}

	err = CreateJSONResponse(jsonGroups, w, http.StatusOK)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	log.Printf("User %v fetched groups for feed", userID)
}
