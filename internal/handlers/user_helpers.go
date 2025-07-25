package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/TheJa750/PrayerPals/internal/auth"
	"github.com/TheJa750/PrayerPals/internal/database"
	"github.com/google/uuid"
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
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return uuid.Nil, fmt.Errorf("getUserIDFromToken: %w", err)
	}

	userID, err := auth.ValidateJWT(tokenString, a.JWTSecret)
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
