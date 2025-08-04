package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/TheJa750/PrayerPals/internal/database"
	"github.com/google/uuid"
)

var ErrUnauthorizedDelete = errors.New("user is not authorized to delete this post")
var ErrPostNotFound = errors.New("post not found")

func (a *APIConfig) createPost(ctx context.Context, groupID, userID uuid.UUID, content string) (Post, error) {
	// Validate user in group (future: add role check in helper function)
	isMember, err := a.verifyUserInGroup(ctx, userID, groupID)
	if err != nil {
		return Post{}, err
	}
	if !isMember {
		return Post{}, ErrUserNotMember
	}

	// Create the post in the database
	post, err := a.DBQueries.CreatePost(ctx, database.CreatePostParams{
		GroupID: groupID,
		UserID:  userID,
		Content: content,
	})
	if err != nil {
		return Post{}, err
	}

	// Convert database post to API Post type
	jsonPost := Post{
		ID:        post.ID,
		GroupID:   post.GroupID,
		UserID:    post.UserID,
		Content:   post.Content,
		CreatedAt: post.CreatedAt.Time.Format(time.RFC3339),
	}

	return jsonPost, nil
}

func (a *APIConfig) createComment(ctx context.Context, postID, userID uuid.UUID, content string) (Comment, error) {
	// Validate user in group (future: add role check in helper function)
	post, err := a.DBQueries.GetPostByID(ctx, postID)
	if err != nil {
		return Comment{}, err
	}

	isMember, err := a.verifyUserInGroup(ctx, userID, post.GroupID)
	if err != nil {
		return Comment{}, err
	}
	if !isMember {
		return Comment{}, ErrUserNotMember
	}

	// Validate post in group
	isValidPost, err := a.verifyPostInGroup(ctx, postID, post.GroupID)
	if err != nil {
		return Comment{}, err
	}
	if !isValidPost {
		return Comment{}, ErrPostNotFound
	}

	// Create the comment in the database
	parentID := uuid.NullUUID{
		UUID:  postID,
		Valid: true,
	}

	comment, err := a.DBQueries.CreateComment(ctx, database.CreateCommentParams{
		ParentPostID: parentID,
		GroupID:      post.GroupID,
		UserID:       userID,
		Content:      content,
	})
	if err != nil {
		return Comment{}, err
	}

	// Return the created comment as JSON response
	jsonComment := Comment{
		ID:      comment.ID,
		PostID:  comment.ParentPostID.UUID,
		GroupID: comment.GroupID,
		UserID:  userID,
		Content: comment.Content,
	}

	return jsonComment, nil
}

func (a *APIConfig) verifyUserCanDeletePost(ctx context.Context, userID, postID, groupID uuid.UUID) error {
	// Verify if the user is a member of the group
	isMember, err := a.verifyUserInGroup(ctx, userID, groupID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrUserNotMember
	}

	// Check if the user is the owner of the post or an admin in the group
	post, err := a.DBQueries.GetPostByID(ctx, postID)
	if err != nil {
		return err
	}

	if post.UserID == userID {
		return nil // User is the owner of the post
	}

	// Check if the user is an admin in the group
	if err = a.isAdmin(ctx, userID, groupID); err != nil {
		return err // error will be ErrUserNotAdmin or a DB query error
	}

	return ErrUnauthorizedDelete
}

func (a *APIConfig) getCommentsOnPost(ctx context.Context, postID uuid.UUID) ([]Comment, error) {
	parentID := uuid.NullUUID{
		UUID:  postID,
		Valid: true,
	}
	comments, err := a.DBQueries.GetCommentsByPostID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	// Convert database comments to API Comment type
	jsonComments := make([]Comment, len(comments))
	for i, comment := range comments {
		jsonComments[i] = Comment{
			ID:        comment.ID,
			PostID:    postID,
			UserID:    comment.UserID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt.Time.Format(time.RFC3339),
			Author:    comment.Username.String,
		}
	}

	return jsonComments, nil
}

func (a *APIConfig) getPostWithComments(ctx context.Context, userID, postID uuid.UUID) (Post, error) {
	// Fetch the post
	post, err := a.DBQueries.GetPostByID(ctx, postID)
	if err != nil {
		return Post{}, err
	}

	// Verify if the user is a member of the group
	isMember, err := a.verifyUserInGroup(ctx, userID, post.GroupID)
	if err != nil {
		return Post{}, err
	}
	if !isMember {
		return Post{}, ErrUserNotMember
	}

	// Fetch comments for the post
	comments, err := a.getCommentsOnPost(ctx, postID)
	if err != nil {
		return Post{}, err
	}

	// Convert database post to API Post type
	jsonPost := Post{
		ID:        post.ID,
		GroupID:   post.GroupID,
		UserID:    post.UserID,
		Content:   post.Content,
		CreatedAt: post.CreatedAt.Time.Format(time.RFC3339),
		Author:    post.Username.String,
		Comments:  comments,
	}

	return jsonPost, nil
}
