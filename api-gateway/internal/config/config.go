// api-gateway/internal/config/config.go
package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	Environment       string
	AuthServiceURL    string
	UserServiceURL    string
	ProductServiceURL string
	OrderServiceURL   string
	JWTSecret         string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:              os.Getenv("PORT"),
		Environment:       os.Getenv("ENVIRONMENT"),
		AuthServiceURL:    os.Getenv("AUTH_SERVICE_URL"),
		UserServiceURL:    os.Getenv("USER_SERVICE_URL"),
		ProductServiceURL: os.Getenv("PRODUCT_SERVICE_URL"),
		OrderServiceURL:   os.Getenv("ORDER_SERVICE_URL"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
