package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (a *APIConfig) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract group ID from URL parameters
	groupID, err := parseUUIDPathParam(r, "group_id")
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Parse the request body for post content
	postReq, err := ParseJSON[PostRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if postReq.Content == "" {
		http.Error(w, "Post content is required", http.StatusBadRequest)
		return
	}

	jsonPost, err := a.createPost(r.Context(), groupID, userID, postReq.Content)
	if err != nil {
		if errors.Is(err, ErrUserNotMember) {
			http.Error(w, "User not a member of the group", http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	err = CreateJSONResponse(jsonPost, w, http.StatusCreated)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	log.Printf("User %v created post %v in group %v", userID, jsonPost.ID, jsonPost.GroupID)
}

func (a *APIConfig) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body for comment content
	commentReq, err := ParseJSON[PostRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if commentReq.Content == "" {
		http.Error(w, "Comment content is required", http.StatusBadRequest)
		return
	}

	// Extract group ID and post ID from URL parameters
	groupID, err := parseUUIDPathParam(r, "group_id")
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}
	postID, err := parseUUIDPathParam(r, "post_id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Create the comment in the database
	comment, err := a.createComment(r.Context(), postID, userID, commentReq.Content)
	if err != nil {
		if errors.Is(err, ErrUserNotMember) {
			http.Error(w, "User not a member of the group", http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	err = CreateJSONResponse(comment, w, http.StatusCreated)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	log.Printf("User %v created comment %v on post %v in group %v", userID, comment.ID, postID, groupID)
}

func (a *APIConfig) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the post ID and group ID from URL parameters
	groupID, err := parseUUIDPathParam(r, "group_id")
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}
	postID, err := parseUUIDPathParam(r, "post_id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Verify if the user can delete the post
	err = a.verifyUserCanDeletePost(r.Context(), userID, postID, groupID)
	if err != nil {
		if errors.Is(err, ErrUserNotMember) || errors.Is(err, ErrUnauthorizedDelete) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to verify authority to delete posts", http.StatusInternalServerError)
		return
	}

	// Delete the post from the database
	err = a.DBQueries.DeletePost(r.Context(), postID)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	log.Printf("User %v deleted post %v in group %v", userID, postID, groupID)

	// Delete comments associated with the post
	parentID := uuid.NullUUID{
		UUID:  postID,
		Valid: true,
	}

	err = a.DBQueries.DeleteCommentsFromPost(r.Context(), parentID)
	if err != nil {
		http.Error(w, "Failed to delete comments associated with post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("User %v deleted comments associated with post %v in group %v", userID, postID, groupID)
}

func (a *APIConfig) GetCommentsForPostHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract post ID from query parameters
	postID, err := parseUUIDPathParam(r, "post_id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Fetch comments for the post
	comments, err := a.getCommentsOnPost(r.Context(), postID)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	err = CreateJSONResponse(comments, w, http.StatusOK)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	log.Printf("User %v fetched comments for post %v", userID, postID)
}
