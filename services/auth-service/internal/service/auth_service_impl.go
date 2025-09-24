package service

import (
	"context"
	"time"

	// "github.com/Dzaakk/micro-commerce/proto/customer"

	"github.com/Dzaakk/micro-commerce/services/auth-service/internal/dto"
	"github.com/Dzaakk/micro-commerce/services/auth-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type authServiceImpl struct {
	// customerClient customer.CustomerServiceClient
	authRepo      repository.AuthRepository
	tokenService  TokenService
	tokenExpiry   time.Duration
	refreshExpiry time.Duration
}

// func NewAuthService(authRepo repository.AuthRepository, tokenService TokenService, customerClient customer.CustomerServiceClient) AuthService {
func NewAuthService(authRepo repository.AuthRepository, tokenService TokenService) AuthService {
	return &authServiceImpl{
		// customerClient: customerClient,
		authRepo:      authRepo,
		tokenService:  tokenService,
		tokenExpiry:   1 * time.Hour,
		refreshExpiry: 7 * 24 * time.Hour,
	}
}

// GetUserByID implements AuthService.
func (a *authServiceImpl) GetUserByID(ctx context.Context, userID string) (*dto.BasicUser, error) {
	panic("unimplemented")
}

// Login implements AuthService.
func (a *authServiceImpl) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	panic("unimplemented")
}

// Logout implements AuthService.
func (a *authServiceImpl) Logout(ctx context.Context, accessToken string) error {
	panic("unimplemented")
}

// RefreshToken implements AuthService.
func (a *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	panic("unimplemented")
}

// Register implements AuthService.
func (a *authServiceImpl) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	panic("unimplemented")
}

// RevokeRefreshToken implements AuthService.
func (a *authServiceImpl) RevokeRefreshToken(ctx context.Context, token string) error {
	panic("unimplemented")
}

// ValidateToken implements AuthService.
func (a *authServiceImpl) ValidateToken(ctx context.Context, token string) (*jwt.MapClaims, error) {
	panic("unimplemented")
}
