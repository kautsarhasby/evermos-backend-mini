package models

import "github.com/golang-jwt/jwt/v5"

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

type TokenClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type RedisSession struct {
	UserID       uint   `json:"user_id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
