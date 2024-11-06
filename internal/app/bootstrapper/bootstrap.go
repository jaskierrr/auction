package bootstrapper

import (
	"context"
	"log"
	"main/internal/app"
	"main/internal/config"
	"main/internal/domain/services/repositories"
	"main/internal/infrastructure/database"
	dbRepo "main/internal/infrastructure/database/repositories"
	"net"

	"google.golang.org/grpc"
)

type bootstrapper struct {
	config *config.Config
}

type Bootstrapper interface {
	registerAPIServer(cfg config.Config) error
	RunAPI() error
}

func New() Bootstrapper {
	return &bootstrapper{
		config: config.NewConfig(),
	}
}

func (b *bootstrapper) RunAPI() error {

	err := b.registerAPIServer(*b.config)
	if err != nil {
		log.Fatal("cant start server")
	}

	return nil
}

func (b *bootstrapper) registerAPIServer(cfg config.Config) error {
	lis, err := net.Listen("tcp", cfg.ServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	db := database.NewDB().NewConn(context.Background(), cfg)

	userRepo := dbRepo.NewUserRepo(db)

	userService := repositories.NewUserService(userRepo)

	service := app.NewUserService(userService)

	app.RegisterGRPC(s, service)
	// grpc.RegisterGRPC(s)
	log.Printf("server listening at %v", lis.Addr())

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}
