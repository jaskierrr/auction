package bootstrapper

import (
	"context"
	"log"
	"log/slog"
	app "main/internal/app/usercases"
	"main/internal/config"
	service "main/internal/domain/services"
	"main/internal/infrastructure/database"
	auctionRepo "main/internal/infrastructure/database/repositories/auction_repository"
	userRepo "main/internal/infrastructure/database/repositories/user_repository"
	pb "main/pkg/grpc"
	"main/pkg/logger"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type bootstrapper struct {
	config *config.Config
	logger *slog.Logger
	db     database.DB

	user struct {
		repo    userRepo.UserRepo
		service service.UserService
		usecase app.UserUsecase
	}

	auction struct {
		repo    auctionRepo.AuctionRepo
		service service.AuctionService
		usecase app.AuctionUsecase
	}
}

type Bootstrapper interface {
	registerAPIServer(cfg config.Config) error
	RunAPI() error
}

func New() Bootstrapper {
	return &bootstrapper{
		config: config.NewConfig(),
		logger: logger.NewLogger(),
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


	b.db = database.NewDB().NewConn(context.Background(), cfg)

	b.user.repo = userRepo.NewUserRepo(b.db, b.logger)
	b.user.service = service.NewUserService(b.user.repo)
	b.user.usecase = app.NewUserUsecase(b.user.service, b.logger)

	b.auction.repo = auctionRepo.NewAuctionRepo(b.db, b.logger)
	b.auction.service = service.NewAuctionService(b.auction.repo)
	b.auction.usecase = app.NewAuctionUsecase(b.auction.service, b.logger)

	//** start gRPC-Gateway
	go func() {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		 mux := runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(CustomMatcher),
		 )
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		err = pb.RegisterAuctionServiceHandlerFromEndpoint(ctx, mux, cfg.ServerPort, opts)
		if err != nil {
			log.Fatalf("failed to start HTTP gateway: %v", err)
		}

		log.Println("HTTP gateway listening on :8081")
		if err := http.ListenAndServe(":8081", mux); err != nil {
			log.Fatalf("failed to serve HTTP gateway: %v", err)
		}
	}()

	//** start gRPC server
	s := grpc.NewServer()
	app.RegisterGRPC(s, b.user.usecase, b.auction.usecase)
	log.Printf("server listening at %v", lis.Addr())

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}


	return nil
}

func CustomMatcher(key string) (string, bool) {
		switch key {
		case "X-User-Id":
		return key, true
		default:
		return runtime.DefaultHeaderMatcher(key)
	}
 }
