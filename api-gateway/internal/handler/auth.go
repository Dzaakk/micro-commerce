package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/Dzaakk/micro-commerce/services/auth-service/proto"
)

type AuthHandler struct {
	client pb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthHandler(authServiceURL string) (*AuthHandler, error) {
	conn, err := grpc.NewClient(authServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	return &AuthHandler{
		client: pb.NewAuthServiceClient(conn),
		conn:   conn,
	}, nil
}

func (h *AuthHandler) Close() {
	if h.conn != nil {
		h.conn.Close()
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email     string `json:"email" binding:"required,email"`
		Username  string `json:"username" binding:"required,min=3,max=20"`
		Password  string `json:"password" binding:"required,min=6"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.Register(ctx, &pb.RegisterRequest{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})

	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": st.Message(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"expires_in":    resp.ExpiresIn,
		"user":          resp.User,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.Login(ctx, &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": st.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"expires_in":    resp.ExpiresIn,
		"user":          resp.User,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.RefreshToken(ctx, &pb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": st.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"expires_in":    resp.ExpiresIn,
	})
}

func (h *AuthHandler) ValidateToken(c *gin.Context) {
	token := extractToken(c.GetHeader("Authorization"))
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "token required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.ValidateToken(ctx, &pb.ValidateTokenRequest{
		Token: token,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   resp.Valid,
		"user_id": resp.UserId,
		"email":   resp.Email,
		"role":    resp.Role,
	})
}

func (h *AuthHandler) ValidateTokenMiddleware(token string) (*pb.ValidateTokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return h.client.ValidateToken(ctx, &pb.ValidateTokenRequest{
		Token: token,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetUser(ctx, &pb.GetUserRequest{
		Id: userID,
	})

	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": st.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, resp.User)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "update profile endpoint",
	})
}

func extractToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
