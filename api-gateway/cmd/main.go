package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dzaakk/micro-commerce/api-gateway/internal/config"
	"github.com/Dzaakk/micro-commerce/api-gateway/internal/handler"
	"github.com/Dzaakk/micro-commerce/api-gateway/internal/middleware"
	"github.com/Dzaakk/micro-commerce/api-gateway/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	conf := config.Load()

	if conf.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(
		middleware.Logger(),
		middleware.CORS(),
		gin.Recovery(),
	)

	authHandler, err := initializeAuthHandler(conf)
	if err != nil {
		log.Fatalf("Fialed to initialize auth handler: %v", err)
	}
	defer authHandler.Close()

	healthHandler := handler.NewHealthHandler()

	router.SetupRoutes(r, authHandler, healthHandler)

	srv := &http.Server{
		Addr:         getServerAddress(conf.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("API Gatyeway starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server : %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdonw:", err)
	}

	log.Println("Server exited successfully")
}

func initializeAuthHandler(conf *config.Config) (*handler.AuthHandler, error) {
	authServiceURL := conf.AuthServiceURL

	if authServiceURL == "" {
		authServiceURL = "localhost:8081"
	}

	return handler.NewAuthHandler(authServiceURL)
}

func getServerAddress(port string) string {
	if port == "" {
		port = "8080"
	}

	if port[0] != ':' {
		port = ":" + port
	}

	return port
}
