package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/TheJa750/PrayerPals/internal/auth"
	"github.com/TheJa750/PrayerPals/internal/database"
	"github.com/TheJa750/PrayerPals/internal/validation"
	"github.com/google/uuid"
)

var (
	ErrNoGroupFound       = errors.New("no group found for the provided invite code")
	ErrUserIsMember       = errors.New("user is already a member of the group")
	ErrUserKickedOrBanned = errors.New("user is kicked or banned from the group")
)

func (a *APIConfig) issueTokens(user database.User, jwtSecret string, activeTime time.Duration, ctx context.Context) (string, string, error) {
	accessToken, err := auth.MakeJWT(user.ID, jwtSecret, activeTime)
	if err != nil {
		return "", "", fmt.Errorf("issueTokens: error creating JWT: %v", err)
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("issueTokens: error creating refresh token: %v", err)
	}

	_, err = a.DBQueries.CreateUserToken(ctx, database.CreateUserTokenParams{
		UserID: user.ID,
		Token:  refreshToken,
	})
	if err != nil {
		return "", "", fmt.Errorf("issueTokens: error storing refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (a *APIConfig) getUserIDFromToken(r *http.Request) (uuid.UUID, error) {
	token, err := r.Cookie("access_token")
	if err != nil {
		return uuid.Nil, fmt.Errorf("getUserIDFromToken: missing access token cookie: %w", err)
	}
	if token == nil {
		return uuid.Nil, fmt.Errorf("getUserIDFromToken: access token cookie is nil")
	}

	userID, err := auth.ValidateJWT(token.Value, a.JWTSecret)
	if err != nil {
		return uuid.Nil, fmt.Errorf("getUserIDFromToken: %w", err)
	}

	return userID, nil
}

func (a *APIConfig) verifyUserInGroup(ctx context.Context, userID, groupID uuid.UUID) (bool, error) {
	members, err := a.DBQueries.GetGroupMembersIDs(ctx, groupID)
	if err != nil {
		return false, fmt.Errorf("verifyUserInGroup: error retrieving group members: %w", err)
	}

	for _, memberID := range members {
		if memberID == userID {
			// Check if the user is banned or kicked from the group
			modStatus, err := a.DBQueries.GetKickBanStatus(ctx, database.GetKickBanStatusParams{
				UserID:  userID,
				GroupID: groupID,
			})
			if err != nil {
				return false, fmt.Errorf("verifyUserInGroup: error checking kick/ban status: %w", err)
			}
			if modStatus.IsBanned || modStatus.IsKicked {
				return false, ErrUserKickedOrBanned
			}

			return true, nil
		}
	}

	return false, nil
}

func (a *APIConfig) verifyPostInGroup(ctx context.Context, postID, groupID uuid.UUID) (bool, error) {
	post, err := a.DBQueries.GetPostByID(ctx, postID)
	if err != nil {
		return false, fmt.Errorf("verifyPostInGroup: error retrieving post: %w", err)
	}

	if post.GroupID != groupID {
		return false, nil
	}

	return true, nil
}

func (a *APIConfig) joinGroup(ctx context.Context, userID uuid.UUID, role, inviteCode string) (UserJoinGroup, error) {
	// Fetch group by invite code
	group, err := a.getGroupByInviteCode(ctx, inviteCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserJoinGroup{}, ErrNoGroupFound
		}
		return UserJoinGroup{}, fmt.Errorf("joinGroup: error retrieving group by invite code: %w", err)
	}

	// Verify user is not already a member (checked with kick/ban status below)
	isMember := false
	members, err := a.DBQueries.GetGroupMembersIDs(ctx, group.ID)
	if err != nil {
		return UserJoinGroup{}, fmt.Errorf("verifyUserInGroup: error retrieving group members: %w", err)
	}

	for _, memberID := range members {
		if memberID == userID {
			isMember = true
			break
		}
	}

	if isMember {
		// Check if the user is banned or kicked from the group
		modStatus, err := a.DBQueries.GetKickBanStatus(ctx, database.GetKickBanStatusParams{
			UserID:  userID,
			GroupID: group.ID,
		})
		if err != nil {
			return UserJoinGroup{}, fmt.Errorf("joinGroup: error checking kick/ban status: %w", err)
		}

		if modStatus.IsBanned {
			return UserJoinGroup{}, fmt.Errorf("joinGroup: user is banned from the group")
		}

		if modStatus.IsKicked {
			if modStatus.KickedUntil.Time.After(time.Now()) { // Kick is still active
				return UserJoinGroup{}, fmt.Errorf("joinGroup: user is kicked from the group until %v", modStatus.KickedUntil.Time)
			} else { // Kick has expired, reset kick status
				err = a.DBQueries.ResetKickStatus(ctx, database.ResetKickStatusParams{
					UserID:  userID,
					GroupID: group.ID,
				})
				if err != nil {
					return UserJoinGroup{}, fmt.Errorf("joinGroup: error resetting kick status: %w", err)
				}
			}
		}

		return UserJoinGroup{}, ErrUserIsMember
	}

	// Add user to the group
	err = a.DBQueries.AddUserToGroup(ctx, database.AddUserToGroupParams{
		UserID:  userID,
		GroupID: group.ID,
		Role:    role,
	})
	if err != nil {
		return UserJoinGroup{}, fmt.Errorf("joinGroup: error adding user to group: %w", err)
	}

	jsonResponse := UserJoinGroup{
		UserID:    userID,
		GroupID:   group.ID,
		GroupName: group.Name,
		Role:      role,
	}

	return jsonResponse, nil
}

func (a *APIConfig) createUser(ctx context.Context, req UserRequest) (User, error) {
	// Hash the password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return User{}, fmt.Errorf("createUser: error hashing password: %w", err)
	}

	// Add user to the database
	userData, err := a.DBQueries.CreateUser(ctx, database.CreateUserParams{
		Username:       req.Username,
		Email:          strings.ToLower(req.Email),
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return User{}, fmt.Errorf("createUser: error creating user: %w", err)
	}

	// Create JSON response
	jsonUser := User{
		ID:       userData.ID,
		Username: userData.Username,
		Email:    userData.Email,
	}

	return jsonUser, nil
}

func (a *APIConfig) updateUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	// Validate new password
	validation := validation.ValidatePassword(newPassword)
	if !validation.IsValid {
		return fmt.Errorf("updateUserPassword: invalid password: %s", strings.Join(validation.Errors, ", "))
	}

	// Hash the new password
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("updateUserPassword: error hashing new password: %w", err)
	}

	// Update the user's password in the database
	err = a.DBQueries.UpdateUserPassword(ctx, database.UpdateUserPasswordParams{
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		return fmt.Errorf("updateUserPassword: error updating user password: %w", err)
	}

	return nil
}

func (a *APIConfig) updateUsername(ctx context.Context, userID uuid.UUID, newUsername string) error {
	// Validate new username
	validation := validation.ValidateUsername(newUsername)
	if !validation.IsValid {
		return fmt.Errorf("updateUsername: invalid username: %s", strings.Join(validation.Errors, ", "))
	}

	// Update the user's username in the database
	err := a.DBQueries.UpdateUsername(ctx, database.UpdateUsernameParams{
		Username: newUsername,
		ID:       userID,
	})
	if err != nil {
		return fmt.Errorf("updateUsername: error updating username: %w", err)
	}

	return nil
}

func (a *APIConfig) fetchGroupsForUser(ctx context.Context, userID uuid.UUID) ([]Group, error) {
	groups, err := a.DBQueries.GetGroupsForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetchGroupsForUser: error retrieving groups for user: %w", err)
	}

	var jsonGroups []Group
	for _, group := range groups {
		jsonGroups = append(jsonGroups, Group{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description.String,
			OwnerID:     group.OwnerID.UUID,
		})
	}

	return jsonGroups, nil
}
