package grpc

import (
	"google.golang.org/grpc"
)

func RegisterGRPC(grpc *grpc.Server) {
	RegisterAuctionServiceServer(grpc, UnimplementedAuctionServiceServer{})
}
