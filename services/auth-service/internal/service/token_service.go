package service

import (
	"errors"
	"time"

	"github.com/Dzaakk/micro-commerce/services/auth-service/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateTokens(user *model.User) (*model.TokenResponse, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
	ValidateRefreshToken(tokenString string) (jwt.MapClaims, error)
}

type tokenService struct {
	jwtSecret []byte
}

func NewTokenService(secret string) TokenService {
	return &tokenService{
		jwtSecret: []byte(secret),
	}
}

func (s *tokenService) GenerateTokens(user *model.User) (*model.TokenResponse, error) {
	// Access token (expires in 15 minutes)
	accessClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
		"type":     "access",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Refresh token (expires in 7 days)
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &model.TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    900, // 15 minutes in seconds
	}, nil
}

func (s *tokenService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check if it's an access token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

func (s *tokenService) ValidateRefreshToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check if it's a refresh token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
