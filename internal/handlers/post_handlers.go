package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/TheJa750/PrayerPals/internal/database"
	"github.com/google/uuid"
)

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
	jsonPost := Post{
		ID:      post.ID,
		GroupID: post.GroupID,
		UserID:  post.UserID,
	}

	err = CreateJSONResponse(jsonPost, w, http.StatusCreated)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	log.Printf("User %v created post %v in group %v", userID, post.ID, postReq.GroupID)
}

func (a *APIConfig) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body for comment content
	commentReq, err := ParseJSON[CommentRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if commentReq.PostID == uuid.Nil {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	// Validate user in group (future: add role check in helper function)
	isMember, err := a.verifyUserInGroup(r.Context(), userID, commentReq.GroupID)
	if err != nil {
		http.Error(w, "Failed to verify group membership", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "User not a member of the group", http.StatusForbidden)
		return
	}

	// Validate post in group
	isValidPost, err := a.verifyPostInGroup(r.Context(), commentReq.PostID, commentReq.GroupID)
	if err != nil {
		http.Error(w, "Failed to verify post belongs to group", http.StatusInternalServerError)
		return
	}
	if !isValidPost {
		http.Error(w, "Post does not exist in group", http.StatusForbidden)
		return
	}

	// Create the comment in the database
	parentID := uuid.NullUUID{
		UUID:  commentReq.PostID,
		Valid: true,
	}

	comment, err := a.DBQueries.CreateComment(r.Context(), database.CreateCommentParams{
		ParentPostID: parentID,
		GroupID:      commentReq.GroupID,
		UserID:       userID,
		Content:      commentReq.Content,
	})
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	// Return the created comment as JSON response
	jsonComment := Comment{
		ID:      comment.ID,
		PostID:  comment.ParentPostID.UUID,
		GroupID: commentReq.GroupID,
		UserID:  userID,
	}

	err = CreateJSONResponse(jsonComment, w, http.StatusCreated)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	log.Printf("User %v created comment %v on post %v in group %v", userID, comment.ID, commentReq.PostID, commentReq.GroupID)
}

func (a *APIConfig) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	// Validate JWT and extract user ID
	userID, err := a.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body for post ID
	deleteReq, err := ParseJSON[DeletePostRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify if the user can delete the post
	err = a.verifyUserCanDeletePost(r.Context(), userID, deleteReq.PostID, deleteReq.GroupID)
	if err != nil {
		if errors.Is(err, ErrUserNotMember) || errors.Is(err, ErrUnauthorizedDelete) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to verify authority to delete posts", http.StatusInternalServerError)
		return
	}

	// Delete the post from the database
	err = a.DBQueries.DeletePost(r.Context(), deleteReq.PostID)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	log.Printf("User %v deleted post %v in group %v", userID, deleteReq.PostID, deleteReq.GroupID)

	// Delete comments associated with the post
	parentID := uuid.NullUUID{
		UUID:  deleteReq.PostID,
		Valid: true,
	}

	err = a.DBQueries.DeleteCommentsFromPost(r.Context(), parentID)
	if err != nil {
		http.Error(w, "Failed to delete comments associated with post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("User %v deleted comments associated with post %v in group %v", userID, deleteReq.PostID, deleteReq.GroupID)
}
