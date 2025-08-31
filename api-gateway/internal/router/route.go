package router

import (
	"github.com/Dzaakk/micro-commerce/api-gateway/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, authHandler *handler.AuthHandler, healthHandler *handler.HealthHandler) {

	r.GET("/", healthHandler.Index)
	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/validate", authHandler.ValidateToken)
		}
	}

}
