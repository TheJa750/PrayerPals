package handlers

import (
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

func (a *APIConfig) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body for post content
	postReq, err := ParseJSON[PostRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate user in group (future: add role check in helper function)
	isMember, err := a.verifyUserInGroup(r.Context(), userID, postReq.GroupID)
	if err != nil {
		http.Error(w, "Failed to verify group membership", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "User not a member of the group", http.StatusForbidden)
		return
	}

	// Create the post in the database
	post, err := a.DBQueries.CreatePost(r.Context(), database.CreatePostParams{
		GroupID: postReq.GroupID,
		UserID:  userID,
		Content: postReq.Content,
	})
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// Return the created post as JSON response
	err = CreateJSONResponse(post, w, http.StatusCreated)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	log.Printf("User %v created post %v in group %v", userID, post.ID, postReq.GroupID)
}
