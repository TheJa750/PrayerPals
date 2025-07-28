package handlers

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strconv"

	"github.com/TheJa750/PrayerPals/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var ErrUserNotMember = errors.New("user is not a member of the group")
var ErrUserIsOnlyAdmin = errors.New("cannot leave group as the only admin")
var ErrUserIsLastMember = errors.New("cannot leave group as the last member")
var ErrInvalidRole = errors.New("invalid role specified")
var ErrUserHasRole = errors.New("user already has the specified role")
var ErrUserNotAdmin = errors.New("user is not an admin of group")
var ErrInvalidID = errors.New("missing or invalid UUID parameter")
var ErrCannotModAdmin = errors.New("cannot moderate an admin")

func (a *APIConfig) leaveGroupChecks(ctx context.Context, userID, groupID uuid.UUID) error {
	// Verify if the user is a member of the group
	isMember, err := a.verifyUserInGroup(ctx, userID, groupID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrUserNotMember
	}

	// Check if the user is the owner or admin of the group
	specialRoleUsers, err := a.DBQueries.GetGroupSpecialRoles(ctx, groupID)
	if err != nil {
		return err
	}

	// Check if the user is the only member with a special role
	if len(specialRoleUsers) == 1 && specialRoleUsers[0].UserID == userID {
		return ErrUserIsOnlyAdmin
	}

	// Check if the user is the last member of the group
	members, err := a.DBQueries.GetGroupMembersIDs(ctx, groupID)
	if err != nil {
		return err
	}
	if len(members) == 1 {
		return ErrUserIsLastMember
	}

	return nil
}

func (a *APIConfig) promoteUserChecks(ctx context.Context, userID, groupID uuid.UUID, role string) error {
	// Verify if the user is a member of the group
	isMember, err := a.verifyUserInGroup(ctx, userID, groupID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrUserNotMember
	}

	// Verify the role is a valid role
	roles := getValidRoles()
	if !slices.Contains(roles, role) {
		return ErrInvalidRole
	}

	// Check if the user is already an admin
	userRole, err := a.DBQueries.GetUserGroupRole(ctx, database.GetUserGroupRoleParams{
		UserID:  userID,
		GroupID: groupID,
	})
	if err != nil {
		return err
	}
	if userRole == role {
		return ErrUserHasRole
	}

	return nil
}

func (a *APIConfig) getPostFeed(ctx context.Context, userID, groupID uuid.UUID, limit, offset int) ([]Post, error) {
	// Verify if the user is a member of the group
	isMember, err := a.verifyUserInGroup(ctx, userID, groupID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrUserNotMember
	}

	// Fetch posts for the group
	posts, err := a.DBQueries.GetPostsForFeed(ctx, database.GetPostsForFeedParams{
		GroupID: groupID,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	// Convert database posts to API Post structs
	jsonPosts := make([]Post, len(posts))
	for i, post := range posts {
		jsonPosts[i] = Post{
			ID:      post.ID,
			GroupID: post.GroupID,
			UserID:  post.UserID,
			Content: post.Content,
		}
	}

	return jsonPosts, nil
}

func (a *APIConfig) isAdmin(ctx context.Context, userID, groupID uuid.UUID) error {
	// Check if the user is an admin of the group
	userRole, err := a.DBQueries.GetUserGroupRole(ctx, database.GetUserGroupRoleParams{
		UserID:  userID,
		GroupID: groupID,
	})
	if err != nil {
		return err
	}

	if userRole != "admin" {
		return ErrUserNotAdmin
	}

	return nil
}

func parseIntQueryParam(r *http.Request, key string, defaultVal int) (int, error) {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultVal, nil
	}
	i, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal, err
	}
	return i, nil
}

func parseUUIDPathParam(r *http.Request, key string) (uuid.UUID, error) {
	valStr := mux.Vars(r)[key]
	if valStr == "" {
		return uuid.Nil, ErrInvalidID
	}
	id, err := uuid.Parse(valStr)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (a *APIConfig) moderateUser(ctx context.Context, groupID, targetID, adminID uuid.UUID, action, reason string) error {
	// Check if the admin is an admin of the group
	if err := a.isAdmin(ctx, adminID, groupID); err != nil {
		return err
	}

	// Verify if the target user is a member of the group
	isMember, err := a.verifyUserInGroup(ctx, targetID, groupID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrUserNotMember
	}

	// Verify the target user is not an admin
	targetRole, err := a.DBQueries.GetUserGroupRole(ctx, database.GetUserGroupRoleParams{
		UserID:  targetID,
		GroupID: groupID,
	})
	if err != nil {
		return err
	}
	if targetRole == "admin" {
		return ErrCannotModAdmin
	}

	// Perform the moderation action
	switch action {
	case "kick": // Kick has a length of 7 days
		return a.DBQueries.KickUser(ctx, database.KickUserParams{
			GroupID:      groupID,
			UserID:       targetID,
			ModdedReason: reason,
			ModdedBy:     uuid.NullUUID{UUID: adminID, Valid: true},
		})
	case "ban": // Ban is permanent
		return a.DBQueries.BanUser(ctx, database.BanUserParams{
			GroupID:      groupID,
			UserID:       targetID,
			ModdedReason: reason,
			ModdedBy:     uuid.NullUUID{UUID: adminID, Valid: true},
		})
	default:
		return errors.New("invalid moderation action")
	}
}

func (a *APIConfig) promoteUser(ctx context.Context, groupID, userID uuid.UUID, role string) error {
	// Perform checks before promoting user
	if err := a.promoteUserChecks(ctx, userID, groupID, role); err != nil {
		return err
	}

	// Update user role in the database
	return a.DBQueries.AdjustUserGroupRole(ctx, database.AdjustUserGroupRoleParams{
		UserID:  userID,
		GroupID: groupID,
		Role:    role,
	})
}
