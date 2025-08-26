package service

import (
	"context"
	"errors"
	"time"

	"github.com/Dzaakk/micro-commerce/services/auth-service/internal/model"
	pb "github.com/Dzaakk/micro-commerce/services/auth-service/internal/proto"
	"github.com/Dzaakk/micro-commerce/services/auth-service/internal/repository"

	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error)
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*pb.AuthResponse, error)
	GetUser(ctx context.Context, userID int64) (*pb.User, error)
}

type authService struct {
	repo         repository.UserRepository
	tokenService TokenService
	client       client.Client
}

func NewAuthService(repo repository.UserRepository, tokenService TokenService, client client.Client) AuthService {
	return &authService{
		repo:         repo,
		tokenService: tokenService,
		client:       client,
	}
}

func (s *authService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" || req.Username == "" {
		return nil, errors.New("email, username, and password are required")
	}

	// Check if user exists
	existingUser, _ := s.repo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	existingUser, _ = s.repo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &model.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         "customer",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		logger.Errorf("Failed to create user: %v", err)
		return nil, errors.New("failed to create user")
	}

	// Generate tokens
	tokens, err := s.tokenService.GenerateTokens(user)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	return &pb.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User:         s.userToProto(user),
	}, nil
}

func (s *authService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	// Get user
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Generate tokens
	tokens, err := s.tokenService.GenerateTokens(user)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	return &pb.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User:         s.userToProto(user),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*pb.AuthResponse, error) {
	// Validate refresh token
	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Get user
	user, err := s.repo.GetByID(ctx, int64(userID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Generate new tokens
	tokens, err := s.tokenService.GenerateTokens(user)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	return &pb.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User:         s.userToProto(user),
	}, nil
}

func (s *authService) GetUser(ctx context.Context, userID int64) (*pb.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.userToProto(user), nil
}

func (s *authService) userToProto(user *model.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
