package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
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

	// Validate JWT and get user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Error getting user ID from token: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create JSON response
	group, err := a.createGroup(r.Context(), userID, groupReq)
	if err != nil {
		log.Printf("Error creating group: %v", err)
		http.Error(w, "Error creating group", http.StatusInternalServerError)
		return
	}

	// Add user to the group as admin
	_, err = a.joinGroup(r.Context(), userID, "admin", group.InviteCode)
	if err != nil {
		log.Printf("Error adding user to group: %v", err)

		// if adding user fails, we delete group because no way to access it
		if deleteErr := a.deleteGroup(r.Context(), group.ID); deleteErr != nil {
			log.Printf("Error deleting group after failed join: %v", deleteErr)
			http.Error(w, "Error deleting group after failed join", http.StatusInternalServerError)
			return
		}

		http.Error(w, "Error adding user to group, group deleted", http.StatusInternalServerError)
		return
	}

	if err := CreateJSONResponse(group, w, http.StatusCreated); err != nil {
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

	// Get group ID and target user ID from URL path
	groupID, err := parseUUIDPathParam(r, "group_id")
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}
	targetUserID, err := parseUUIDPathParam(r, "user_id")
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Check if the user is admin in group
	if err = a.isAdmin(r.Context(), userID, groupID); err != nil {
		if errors.Is(err, ErrUserNotAdmin) {
			http.Error(w, "User is not an admin of the group", http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to check admin status", http.StatusInternalServerError)
		return
	}

	// Perform checks and promote user
	role := strings.ToLower(promoteReq.Role)
	err = a.promoteUser(r.Context(), groupID, targetUserID, role)
	if err != nil {
		if errors.Is(err, ErrUserNotMember) {
			http.Error(w, "Target user is not a member of the group", http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to promote user", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusNoContent)
	log.Printf("User %v promoted to %s in group %v", targetUserID, role, groupID)

}

func (a *APIConfig) GetPostFeedHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the group ID from the URL path
	groupID, err := parseUUIDPathParam(r, "group_id")
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Parse query parameters for pagination
	limit, err := parseIntQueryParam(r, "limit", 10)
	if err != nil {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}
	offset, err := parseIntQueryParam(r, "offset", 0)
	if err != nil {
		http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
		return
	}

	// Perform checks and get posts for the group
	posts, err := a.getPostFeed(r.Context(), userID, groupID, limit, offset)
	if err != nil {
		if errors.Is(err, ErrUserNotMember) {
			http.Error(w, "User is not a member of the group", http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to get post feed", http.StatusInternalServerError)
		return
	}

	// Create JSON response
	if err := CreateJSONResponse(posts, w, http.StatusOK); err != nil {
		log.Printf("Error creating JSON response: %v", err)
		return
	}
}

func (a *APIConfig) GetPostCountHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the group ID from the URL path
	groupID, err := parseUUIDPathParam(r, "group_id")
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Get the post count for the group
	count, err := a.getGroupPostCount(r.Context(), userID, groupID)
	if err != nil {
		if errors.Is(err, ErrUserNotMember) {
			http.Error(w, "User is not a member of the group", http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to get post count", http.StatusInternalServerError)
		return
	}

	// Create JSON response
	response := PostCountResponse{
		PostCount: count,
	}
	if err := CreateJSONResponse(response, w, http.StatusOK); err != nil {
		log.Printf("Error creating JSON response: %v", err)
		return
	}
	log.Printf("Post count retrieved successfully for group %v by user %v", groupID, userID)
}

func (a *APIConfig) DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the group ID from the URL path
	groupID, err := parseUUIDPathParam(r, "group_id")
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Check if the user is an admin of the group
	if err = a.isAdmin(r.Context(), userID, groupID); err != nil {
		if errors.Is(err, ErrUserNotAdmin) {
			http.Error(w, "User is not an admin of the group", http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to check admin status", http.StatusInternalServerError)
		return
	}

	// Delete the group from the database
	err = a.DBQueries.DeleteGroup(r.Context(), groupID)
	if err != nil {
		http.Error(w, "Failed to delete group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("Group %v deleted successfully by user %v", groupID, userID)
}

func (a *APIConfig) ModerateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body for moderation details
	moderateReq, err := ParseJSON[ModerateUserRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get group/target IDs from URL path
	groupID, err := parseUUIDPathParam(r, "group_id")
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	targetUserID, err := parseUUIDPathParam(r, "user_id")
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Moderate the user
	err = a.moderateUser(
		r.Context(),
		groupID,
		targetUserID,
		userID,
		strings.ToLower(moderateReq.Action),
		moderateReq.Reason,
	)
	if err != nil {
		if errors.Is(err, ErrUserNotAdmin) {
			http.Error(w, "User is not an admin of the group", http.StatusForbidden)
			return
		}
		if errors.Is(err, ErrUserNotMember) {
			http.Error(w, "Target user is not a member of the group", http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to moderate user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("User %v moderated in group %v by admin %v", targetUserID, groupID, userID)
}

func (a *APIConfig) GroupFromInviteCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the invite code from the URL path
	inviteCode, err := parseInviteCodePathParam(r, "invite_code")
	if err != nil {
		http.Error(w, "Invalid invite code", http.StatusBadRequest)
		return
	}

	// Fetch group by invite code
	group, err := a.getGroupByInviteCode(r.Context(), inviteCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Group not found for invite code", http.StatusNotFound)
			return
		}
		http.Error(w, "Error retrieving group by invite code", http.StatusInternalServerError)
		return
	}

	// Create JSON response
	if err := CreateJSONResponse(group, w, http.StatusOK); err != nil {
		log.Printf("Error creating JSON response: %v", err)
		return
	}

	log.Printf("Group retrieved successfully for invite code: %s", inviteCode)
}

func (a *APIConfig) GetGroupInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the group ID from the URL path
	groupID, err := parseUUIDPathParam(r, "group_id")
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Fetch group information
	groupInfo, err := a.getGroupByID(r.Context(), userID, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrUserNotMember) {
			http.Error(w, "User is not a member of the group", http.StatusForbidden)
			return
		}
		http.Error(w, "Error retrieving group information", http.StatusInternalServerError)
		return
	}

	// Create JSON response
	if err := CreateJSONResponse(groupInfo, w, http.StatusOK); err != nil {
		log.Printf("Error creating JSON response: %v", err)
		return
	}

	log.Printf("Group information retrieved successfully for group ID: %s", groupID)
}

func (a *APIConfig) GetGroupMembersHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the group ID from the URL path
	groupID, err := parseUUIDPathParam(r, "group_id")
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Check if the user is a member of the group
	isMember, err := a.verifyUserInGroup(r.Context(), userID, groupID)
	if err != nil {
		http.Error(w, "Error verifying group membership", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "User is not a member of the group", http.StatusForbidden)
		return
	}

	// Fetch group members
	members, err := a.getGroupMembers(r.Context(), groupID)
	if err != nil {
		http.Error(w, "Error retrieving group members", http.StatusInternalServerError)
		return
	}

	// Create JSON response
	if err := CreateJSONResponse(members, w, http.StatusOK); err != nil {
		log.Printf("Error creating JSON response: %v", err)
		return
	}

	log.Printf("Group members retrieved successfully for group ID: %s", groupID)
}

func (a *APIConfig) GetUserGroupRoleHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the group ID and user ID from the URL path
	groupID, err := parseUUIDPathParam(r, "group_id")
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}
	targetUserID, err := parseUUIDPathParam(r, "user_id")
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if userID != targetUserID {
		http.Error(w, "Unauthorized to view other user's role", http.StatusForbidden)
		return
	}

	// Get the user's role in the group
	role, err := a.getUserGroupRole(r.Context(), userID, groupID)
	if err != nil {
		if errors.Is(err, ErrUserNotMember) {
			http.Error(w, "User is not a member of the group", http.StatusForbidden)
			return
		}
		http.Error(w, "Error retrieving user role in group", http.StatusInternalServerError)
		return
	}

	response := GroupMember{
		UserID: userID,
		Role:   role,
	}

	// Create JSON response
	if err := CreateJSONResponse(response, w, http.StatusOK); err != nil {
		log.Printf("Error creating JSON response: %v", err)
		return
	}
}

func (a *APIConfig) ChangeInviteCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the group ID from the URL path
	groupID, err := parseUUIDPathParam(r, "group_id")
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Parse the new invite code from the request body
	code, err := ParseJSON[UpdateInviteCodeRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = a.updateGroupInviteCode(r.Context(), userID, groupID, code.InviteCode)
	if err != nil {
		if errors.Is(err, ErrUserNotAdmin) {
			http.Error(w, "User is not an admin of the group", http.StatusForbidden)
			return
		}
		http.Error(w, "Error updating invite code", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("Invite code updated successfully for group %v by user %v", groupID, userID)
}

func (a *APIConfig) ChangeGroupRulesHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the group ID from the URL path
	groupID, err := parseUUIDPathParam(r, "group_id")
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Parse the new rules from the request body
	rules, err := ParseJSON[UpdateGroupRulesRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the group rules in the database
	err = a.changeGroupRules(r.Context(), userID, groupID, rules.Rules)
	if err != nil {
		if errors.Is(err, ErrUserNotAdmin) {
			http.Error(w, "User is not an admin of the group", http.StatusForbidden)
			return
		}
		if errors.Is(err, ErrRulesTooLong) {
			http.Error(w, "Rules text is too long", http.StatusBadRequest)
			return
		}
		http.Error(w, "Error updating group rules", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("Group rules updated successfully for group %v by user %v", groupID, userID)
}

func (a *APIConfig) ChangeGroupDescriptionHandler(w http.ResponseWriter, r *http.Request) {}
