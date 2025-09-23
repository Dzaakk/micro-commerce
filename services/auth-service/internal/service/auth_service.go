package service

import (
	"context"

	"github.com/Dzaakk/micro-commerce/services/auth-service/internal/dto"
	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles user authentication
type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
	Logout(ctx context.Context, accessToken string) error
	ValidateToken(ctx context.Context, token string) (*jwt.MapClaims, error)
	GetUserByID(ctx context.Context, userID string) (*dto.BasicUser, error)
	RevokeRefreshToken(ctx context.Context, token string) error
}

// OAuthService handles OAuth 2.0 operations
type OAuthService interface {
	// Authorization Code Flow
	GenerateAuthorizationCode(ctx context.Context, clientID, userID, redirectURI, scope, state string) (string, error)
	ExchangeCodeForToken(ctx context.Context, req *dto.OAuthTokenRequest) (*dto.OAuthTokenResponse, error)

	// Token Management
	ValidateAccessToken(ctx context.Context, tokenString string) (*jwt.MapClaims, error)
	RevokeToken(ctx context.Context, token string) error
	IntrospectToken(ctx context.Context, token string) (map[string]interface{}, error)

	// Client Management
	ValidateClient(ctx context.Context, clientID string) error
	ValidateRedirectURI(ctx context.Context, clientID, redirectURI string) error
}

// TokenService handles JWT operations
type TokenService interface {
	GenerateAccessToken(userID, clientID, scope string) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateToken(token string) (*jwt.MapClaims, error)
	ExtractClaims(tokenString string) (*jwt.MapClaims, error)
}

// UserService handles user operations
type UserService interface {
	CreateUser(ctx context.Context, req *dto.RegisterRequest) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*dto.BasicUser, error)
	GetUserByID(ctx context.Context, id string) (*dto.BasicUser, error)
	UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error
	ValidatePassword(ctx context.Context, email, password string) (string, error)
}

// SessionService handles session management
type SessionService interface {
	CreateSession(ctx context.Context, userID string) (string, error)
	GetSession(ctx context.Context, sessionID string) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error
	ExtendSession(ctx context.Context, sessionID string) error
}
