package utils

import (
	"context"
	"evermos-hasby/ecommerce/internal/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var TokenType = map[string]time.Duration{
	"token":        time.Minute * 15,
	"refreshToken": (time.Hour * 24) * 365,
}

func GenerateToken(ctx context.Context, userID uint, username, role, tokenType string) (string, error) {
	return GenerateTokenWithID(ctx, userID, username, role, tokenType, uuid.NewString())
}

func GenerateTokenWithID(ctx context.Context, userID uint, username, role, tokenType, tokenID string) (string, error) {
	secretKey := []byte(GetEnv("APP_SECRET", ""))
	claims := models.TokenClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			Issuer:    GetEnv("APP_NAME", ""),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenType[tokenType])),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	resultToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error %v", err)
	}

	return resultToken, nil
}

func ValidateToken(ctx context.Context, token string) (*models.TokenClaims, error) {
	secretKey := []byte(GetEnv("APP_SECRET", ""))
	jwtToken, err := jwt.ParseWithClaims(token, &models.TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("failed to validate method JWT %v", t.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to parse %v", err)
	}

	claimToken, ok := jwtToken.Claims.(*models.TokenClaims)
	if !ok || !jwtToken.Valid {
		return nil, fmt.Errorf("Token Invalid")
	}

	return claimToken, nil
}
