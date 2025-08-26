package handler

import (
	"context"
	"encoding/json"

	pb "github.com/Dzaakk/micro-commerce/services/auth-service/internal/proto"
	"github.com/Dzaakk/micro-commerce/services/auth-service/internal/service"

	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/logger"
)

type AuthHandler struct {
	authService  service.AuthService
	tokenService service.TokenService
	broker       broker.Broker
}

func NewAuthHandler(authService service.AuthService, tokenService service.TokenService, broker broker.Broker) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		tokenService: tokenService,
		broker:       broker,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest, rsp *pb.AuthResponse) error {
	logger.Info("Register request received")

	// Call service
	result, err := h.authService.Register(ctx, req)
	if err != nil {
		logger.Errorf("Registration failed: %v", err)
		return err
	}

	// Map response
	rsp.AccessToken = result.AccessToken
	rsp.RefreshToken = result.RefreshToken
	rsp.ExpiresIn = result.ExpiresIn
	rsp.User = result.User

	// Publish event to RabbitMQ
	h.publishEvent("user.registered", map[string]interface{}{
		"user_id": result.User.Id,
		"email":   result.User.Email,
	})

	return nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.AuthResponse) error {
	logger.Info("Login request received")

	result, err := h.authService.Login(ctx, req)
	if err != nil {
		logger.Errorf("Login failed: %v", err)
		return err
	}

	rsp.AccessToken = result.AccessToken
	rsp.RefreshToken = result.RefreshToken
	rsp.ExpiresIn = result.ExpiresIn
	rsp.User = result.User

	// Publish event
	h.publishEvent("user.logged_in", map[string]interface{}{
		"user_id": result.User.Id,
		"email":   result.User.Email,
	})

	return nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest, rsp *pb.ValidateTokenResponse) error {
	logger.Info("ValidateToken request received")

	claims, err := h.tokenService.ValidateToken(req.Token)
	if err != nil {
		rsp.Valid = false
		return nil
	}

	userID, _ := claims["user_id"].(float64)
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	rsp.Valid = true
	rsp.UserId = int64(userID)
	rsp.Email = email
	rsp.Role = role

	return nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest, rsp *pb.AuthResponse) error {
	logger.Info("RefreshToken request received")

	result, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		logger.Errorf("Token refresh failed: %v", err)
		return err
	}

	rsp.AccessToken = result.AccessToken
	rsp.RefreshToken = result.RefreshToken
	rsp.ExpiresIn = result.ExpiresIn
	rsp.User = result.User

	return nil
}

func (h *AuthHandler) GetUser(ctx context.Context, req *pb.GetUserRequest, rsp *pb.UserResponse) error {
	logger.Info("GetUser request received")

	user, err := h.authService.GetUser(ctx, req.Id)
	if err != nil {
		logger.Errorf("Get user failed: %v", err)
		return err
	}

	rsp.User = user

	return nil
}

func (h *AuthHandler) publishEvent(topic string, data interface{}) {
	body, _ := json.Marshal(data)
	msg := &broker.Message{
		Header: map[string]string{
			"event": topic,
		},
		Body: body,
	}

	if err := h.broker.Publish(topic, msg); err != nil {
		logger.Errorf("Failed to publish event %s: %v", topic, err)
	} else {
		logger.Infof("Event published: %s", topic)
	}
}
