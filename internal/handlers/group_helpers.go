package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"

	"github.com/TheJa750/PrayerPals/internal/database"
	"github.com/TheJa750/PrayerPals/internal/validation"
	"github.com/google/uuid"
)

var (
	ErrUserNotMember    = errors.New("user is not a member of the group")
	ErrUserIsOnlyAdmin  = errors.New("cannot leave group as the only admin")
	ErrUserIsLastMember = errors.New("cannot leave group as the last member")
	ErrInvalidRole      = errors.New("invalid role specified")
	ErrUserHasRole      = errors.New("user already has the specified role")
	ErrUserNotAdmin     = errors.New("user is not an admin of group")
	ErrInvalidID        = errors.New("missing or invalid UUID parameter")
	ErrCannotModAdmin   = errors.New("cannot moderate an admin")
	ErrInvalidJWT       = errors.New("invalid JWT token")
	ErrRulesTooLong     = errors.New("group rules cannot exceed 1500 characters")
)

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
			ID:           post.ID,
			GroupID:      post.GroupID,
			UserID:       post.UserID,
			Content:      post.Content,
			CreatedAt:    post.CreatedAt.Time.Format(time.RFC3339),
			Author:       post.Username.String,
			CommentCount: post.CommentCount,
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

func generateInviteCode(customPrefix string) string {
	// Set some defaults
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	prefix := ""
	randomChars := 3

	// "" means use default prefix and 6 random characters
	if customPrefix == "" {
		prefix = "INV"
		randomChars = 6
	} else {
		if len(customPrefix) > 6 {
			customPrefix = customPrefix[:6] // Truncate to 6 characters
		}
		prefix = strings.ToUpper(customPrefix)
		if len(prefix) < 6 {
			randomChars = 9 - len(prefix) // Ensure total length is 9
		}
	}

	suffix := make([]byte, randomChars)
	for i := range suffix {
		suffix[i] = charset[rand.Intn(len(charset))]
	}

	return prefix + string(suffix)

}

func (a *APIConfig) createGroup(ctx context.Context, userID uuid.UUID, req GroupRequest) (Group, error) {
	// Setting up query parameters
	// Description can be null, so we use sql.NullString
	valid := true
	if req.Description == "" {
		valid = false
	}

	// Create the group in the database
	group, err := a.DBQueries.CreateGroup(ctx, database.CreateGroupParams{
		Name: req.Name,
		Description: sql.NullString{
			String: req.Description,
			Valid:  valid,
		},
		OwnerID: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
		InviteCode: generateInviteCode(""), // Generate a random invite code in the form of "INVxxxxxx"
	})
	if err != nil {
		if strings.Contains(err.Error(), "invite_code") && strings.Contains(err.Error(), "unique") {
			return Group{}, fmt.Errorf("invite code collision, please try again")
		}
		return Group{}, err
	}

	jsonGroup := Group{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description.String,
		OwnerID:     group.OwnerID.UUID,
		InviteCode:  group.InviteCode,
	}

	return jsonGroup, nil
}

func (a *APIConfig) deleteGroup(ctx context.Context, groupID uuid.UUID) error {
	// Assumes the group exists and the user has permission to delete it
	err := a.DBQueries.DeleteGroup(ctx, groupID)
	if err != nil {
		return err
	}

	return nil
}

func (a *APIConfig) getGroupByInviteCode(ctx context.Context, inviteCode string) (Group, error) {
	// Fetch the group by invite code
	group, err := a.DBQueries.GetGroupByInviteCode(ctx, inviteCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Group{}, fmt.Errorf("group not found for invite code: %s", inviteCode)
		}
		return Group{}, fmt.Errorf("error retrieving group by invite code: %w", err)
	}

	jsonGroup := Group{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description.String,
		OwnerID:     group.OwnerID.UUID,
		InviteCode:  group.InviteCode,
	}

	return jsonGroup, nil
}

func (a *APIConfig) updateGroupInviteCode(ctx context.Context, userID, groupID uuid.UUID, newInviteCode string) error {
	// Verify the user is an admin of the group
	err := a.isAdmin(ctx, userID, groupID)
	if err != nil {
		return err
	}

	// Trim whitespace and capitalize the invite code
	newInviteCode = strings.TrimSpace(strings.ToUpper(newInviteCode))

	// Validate the new invite code
	result := validation.ValidateInviteCode(newInviteCode)

	if !result.IsValid {
		return fmt.Errorf("invalid invite code: %s", strings.Join(result.Errors, ", "))
	}

	newInviteCode = generateInviteCode(newInviteCode) // Generate a new invite code with the provided prefix

	// Update the invite code for the group
	err = a.DBQueries.UpdateGroupInviteCode(ctx, database.UpdateGroupInviteCodeParams{
		ID:         groupID,
		InviteCode: newInviteCode,
	})
	if err != nil {
		return fmt.Errorf("error updating group invite code: %w", err)
	}

	return nil
}

func (a *APIConfig) getGroupByID(ctx context.Context, userID, groupID uuid.UUID) (Group, error) {
	// Verify if the user is a member of the group
	isMember, err := a.verifyUserInGroup(ctx, userID, groupID)
	if err != nil {
		return Group{}, err
	}
	if !isMember {
		return Group{}, ErrUserNotMember
	}

	// Fetch the group by ID
	group, err := a.DBQueries.GetGroupByID(ctx, groupID)
	if err != nil {
		return Group{}, fmt.Errorf("error retrieving group by ID: %w", err)
	}

	jsonGroup := Group{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description.String,
		OwnerID:     group.OwnerID.UUID,
		InviteCode:  group.InviteCode,
	}

	return jsonGroup, nil
}

func (a *APIConfig) getGroupMembers(ctx context.Context, groupID uuid.UUID) ([]GroupMember, error) {
	// Fetch group members
	members, err := a.DBQueries.GetActiveMembers(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving group members: %w", err)
	}

	// Convert database members to API GroupMember structs
	jsonMembers := make([]GroupMember, len(members))
	for i, member := range members {
		jsonMembers[i] = GroupMember{
			UserID:   member.ID,
			Role:     member.Role,
			Email:    member.Email,
			Username: member.Username,
		}
	}

	return jsonMembers, nil
}

func (a *APIConfig) getUserGroupRole(ctx context.Context, userID, groupID uuid.UUID) (string, error) {
	// Verify if the user is a member of the group
	isMember, err := a.verifyUserInGroup(ctx, userID, groupID)
	if err != nil {
		return "", err
	}
	if !isMember {
		return "", ErrUserNotMember
	}

	// Fetch the user's role in the group
	role, err := a.DBQueries.GetUserGroupRole(ctx, database.GetUserGroupRoleParams{
		UserID:  userID,
		GroupID: groupID,
	})
	if err != nil {
		return "", fmt.Errorf("error retrieving user role in group: %w", err)
	}

	return role, nil
}

func (a *APIConfig) changeGroupRules(ctx context.Context, userID, groupID uuid.UUID, newRules string) error {
	// Verify the user is an admin of the group
	err := a.isAdmin(ctx, userID, groupID)
	if err != nil {
		return err
	}

	// Validate the new rules
	if len(newRules) > 1500 { // Example validation: max length of 1500 characters
		return ErrRulesTooLong
	}

	err = a.DBQueries.UpdateGroupRules(ctx, database.UpdateGroupRulesParams{
		ID:        groupID,
		RulesInfo: newRules,
	})
	if err != nil {
		return fmt.Errorf("error updating group rules: %w", err)
	}

	return nil
}
