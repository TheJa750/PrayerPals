package handlers

import (
	"context"
	"errors"
	"slices"

	"github.com/TheJa750/PrayerPals/internal/database"
	"github.com/google/uuid"
)

var ErrUserNotMember = errors.New("user is not a member of the group")
var ErrUserIsOnlyAdmin = errors.New("cannot leave group as the only admin")
var ErrUserIsLastMember = errors.New("cannot leave group as the last member")
var ErrInvalidRole = errors.New("invalid role specified")
var ErrUserHasRole = errors.New("user already has the specified role")

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
