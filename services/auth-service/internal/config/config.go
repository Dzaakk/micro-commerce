package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	JWTExpiry     string
	RabbitMQURL   string
	ConsulAddress string
	Environment   string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:          os.Getenv("PORT"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiry:     os.Getenv("JWT_EXPIRY"),
		RabbitMQURL:   os.Getenv("RABBITMQ_URL"),
		ConsulAddress: os.Getenv("CONSUL_ADDRESS"),
		Environment:   os.Getenv("ENVIRONMENT"),
	}
}
