package handlers

import (
	"context"
	"errors"

	"github.com/TheJa750/PrayerPals/internal/database"
	"github.com/google/uuid"
)

var ErrUnauthorizedDelete = errors.New("user is not authorized to delete this post")

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
	userRole, err := a.DBQueries.GetUserGroupRole(ctx, database.GetUserGroupRoleParams{
		UserID:  userID,
		GroupID: groupID,
	})
	if err != nil {
		return err
	}

	if userRole == "admin" {
		return nil // User is an admin in the group
	}

	return ErrUnauthorizedDelete
}
