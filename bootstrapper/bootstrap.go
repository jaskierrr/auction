package bootstrapper

import (
	"context"
	"log"
	"log/slog"
	"main/config"
	"main/internal/handlers"
	"main/internal/infrastructure/database"
	auctionRepo "main/internal/repositories/auction"
	userRepo "main/internal/repositories/user"
	auction_service "main/internal/services/auction"
	user_service "main/internal/services/user"
	pb "main/pkg/grpc"
	"main/pkg/logger"
	"net"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type bootstrapper struct {
	config    *config.Config
	db        database.DB
	logger    *slog.Logger
	validator *validator.Validate

	user struct {
		repo     userRepo.UserRepo
		service  user_service.UserService
		handlers handlers.UserHandlers
	}

	auction struct {
		repo     auctionRepo.AuctionRepo
		service  auction_service.AuctionService
		handlers handlers.AuctionHandlers
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
		return err
	}

	return nil
}

func (b *bootstrapper) registerAPIServer(cfg config.Config) error {
	lis, err := net.Listen("tcp", cfg.ServerPort)
	if err != nil {
		return err
	}

	b.db = database.NewDB().NewConn(context.Background(), cfg)
	b.validator = validator.New(validator.WithRequiredStructEnabled())

	b.user.repo = userRepo.NewUserRepo(b.db, b.logger)
	b.user.service = user_service.NewUserService(b.user.repo, b.logger)
	b.user.handlers = handlers.NewUserHandlers(b.user.service, b.logger, b.validator)

	b.auction.repo = auctionRepo.NewAuctionRepo(b.db, b.logger)
	b.auction.service = auction_service.NewAuctionService(b.auction.repo, b.logger)
	b.auction.handlers = handlers.NewAuctionHandlers(b.auction.service, b.logger, b.validator)

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

		log.Println("HTTP gateway listening on :8081")
		err = http.ListenAndServe(":8081", mux)
	}()
	if err != nil {
		return err
	}

	//** start gRPC server
	s := grpc.NewServer()
	handlers.RegisterGRPC(s, b.user.handlers, b.auction.handlers)
	log.Printf("server listening at %v", lis.Addr())

	err = s.Serve(lis)
	if err != nil {
		return err
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
