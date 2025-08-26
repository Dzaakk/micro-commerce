package main

import (
	"database/sql"

	"github.com/Dzaakk/micro-commerce/services/auth-service/config"
	"github.com/Dzaakk/micro-commerce/services/auth-service/internal/handler"
	"github.com/Dzaakk/micro-commerce/services/auth-service/internal/repository/postgres"
	pb "github.com/Dzaakk/micro-commerce/services/auth-service/proto"
	"github.com/Dzaakk/micro-commerce/services/auth-service/service"
	"go.uber.org/zap"

	"github.com/go-micro/plugins/v4/broker/rabbitmq"
	"github.com/go-micro/plugins/v4/registry/consul"
	_ "github.com/lib/pq"
	"go-micro.dev/v4"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
)

func main() {

	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	consulRegistry := consul.NewRegistry(
		registry.Addrs(cfg.ConsulAddress),
	)

	rabbitBroker := rabbitmq.NewBroker(
		broker.Addrs(cfg.RabbitMQURL),
	)

	srv := micro.NewService(
		micro.Name("auth-service"),
		micro.Version("latest"),
		micro.Registry(consulRegistry),
		micro.Broker(rabbitBroker),
		micro.Address(":"+cfg.Port),
	)

	srv.Init()

	userRepo := postgres.NewUserRepository(db)

	tokenService := service.NewTokenService(cfg.JWTSecret)
	authService := service.NewAuthService(cfg.JWTSecret, tokenService, srv.Client())

	authHandler := handler.NewAuthHandler(authService, tokenService, srv.Server().Options().Broker)

	if err := pb.RegisterAuthServiceServer(srv.Server(), authHandler); err != nil {
		logger.Fatal(err)
	}

	logger.Infof("Starting %s on port %s", srv.Name(), cfg.Port)
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
