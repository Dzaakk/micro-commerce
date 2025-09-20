package repository

import (
	"context"

	"github.com/Dzaakk/micro-commerce/services/auth-service/internal/model"
)

type OAuthRepository interface {
	// OAuth Clients
	GetClientByID(ctx context.Context, clientID string) (*model.OAuthClient, error)
	ValidateClientCredentials(ctx context.Context, clientID, clientSecret string) (*model.OAuthClient, error)

	// Authorization Codes
	CreateAuthorizationCode(ctx context.Context, code *model.AuthorizationCode) error
	GetAuthorizationCode(ctx context.Context, code string) (*model.AuthorizationCode, error)
	DeleteAuthorizationCode(ctx context.Context, code string) error

	// Access Tokens
	CreateAccessToken(ctx context.Context, token *model.AccessToken) error
	GetAccessToken(ctx context.Context, token string) (*model.AccessToken, error)
	DeleteAccessToken(ctx context.Context, token string) error
	DeleteExpiredAccessTokens(ctx context.Context) error

	// Refresh Tokens
	CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRefreshTokenByAccessTokenID(ctx context.Context, accessTokenID string) error
}
