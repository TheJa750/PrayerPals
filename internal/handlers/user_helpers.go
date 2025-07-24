package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/TheJa750/PrayerPals/internal/auth"
	"github.com/TheJa750/PrayerPals/internal/database"
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
