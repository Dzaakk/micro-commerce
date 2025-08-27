package main

import (
	"database/sql"

	"github.com/Dzaakk/micro-commerce/services/auth-service/internal/config"
	"github.com/Dzaakk/micro-commerce/services/auth-service/internal/handler"
	"github.com/Dzaakk/micro-commerce/services/auth-service/internal/repository/postgres"
	"github.com/Dzaakk/micro-commerce/services/auth-service/internal/service"
	pb "github.com/Dzaakk/micro-commerce/services/auth-service/proto"

	_ "github.com/lib/pq"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

func main() {

	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	srv := micro.NewService(
		micro.Name("auth-service"),
		micro.Version("latest"),
		micro.Address(":"+cfg.Port),
	)

	srv.Init()

	userRepo := postgres.NewUserRepository(db)

	tokenService := service.NewTokenService(cfg.JWTSecret)
	authService := service.NewAuthService(userRepo, tokenService, srv.Client())

	authHandler := handler.NewAuthHandler(authService, tokenService, srv.Server().Options().Broker)

	if err := pb.RegisterAuthServiceHandler(srv.Server(), authHandler); err != nil {
		logger.Fatal(err)
	}

	logger.Infof("Starting %s on port %s", srv.Name(), cfg.Port)
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
