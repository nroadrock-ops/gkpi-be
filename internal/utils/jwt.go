package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(userID uint, role string, jemaatID *uint, secret string) (string, string, error) {
	// Access Token (15 menit)
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	}
	if jemaatID != nil {
		accessClaims["jemaat_id"] = *jemaatID
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	// Refresh Token (7 hari)
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	if jemaatID != nil {
		refreshClaims["jemaat_id"] = *jemaatID
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
