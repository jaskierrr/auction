package repositories

import (
	"context"
	"log/slog"
	"main/internal/domain/entities"
	"main/internal/infrastructure/database"
	pb "main/pkg/grpc"
)

const (
	activeLotStatus = "Active"
	closedLotStatus = "Closed"

	activeAuctionStatus = "Active"
	closedAuctionStatus = "Ended"

	createLotQuery = `insert into lots (title, description, starting_bid, seller_id, status)
										values (@title, @description, @starting_bid, @seller_id, @status) returning *`
	getLotQuery = `select * from lots where id = @lot_id`

	createAuctionQuery = `insert into auctions (lot_id, status)
	values (@lot_id, @status) returning id, lot_id, status`

	getBidsQuery = `select * from bids where auction_id = @auction_id`

	endedAuctionQuery = `update auctions set (status, winner_id) = (@status, @winner_id) where id = @auction_id returning *`

	writeTransactionEndedAuctionQuery = `insert into transactions (sender_id, sender_type, recipient_id, recipient_type, amount, transaction_type)
	values (@auction_id, 'Auction', @bidder_id, 'User', @amount, 'Refund')`

	findLotId = `select l.id
	from auctions a
	join lots l on a.lot_id = l.id
	where a.id = @auction_id and a.status = @status`

	closedLotQuery = `update lots set status = 'Closed' where id = $1`

	returnMoney = `update users set balance = balance + @amount where id = @bidder_id returning *`

	findStartingBidQuery = `select l.starting_bid, l.seller_id
	from auctions a
	join lots l on a.lot_id = l.id
	where a.id = @auction_id`

	checkAuctionStatusQuery = `select status from auctions where id = @auction_id`

	editUserBalance = `update users set balance = balance - @amount where id = @bidder_id`

	placeBidQuery = `insert into bids (auction_id, bidder_id, amount)
	values (@auction_id, @bidder_id, @amount)
	on conflict (auction_id, bidder_id)
	do update set amount = bids.amount + EXCLUDED.amount returning *`

	writeTransactionPlaceBidQuery = `insert into transactions (sender_id, sender_type, recipient_id, recipient_type, amount, transaction_type)
	values (@bidder_id, 'User', @auction_id, 'Auction', @amount, 'Payment')`

	getBidQuery = `select * from bids where auction_id = @auction_id`
)

type auctionRepo struct {
	db     database.DB
	logger *slog.Logger
}

type AuctionRepo interface {
	CreateLot(ctx context.Context, in *pb.CreateLotRequest) (entities.Lot, error)
	GetLot(ctx context.Context, in *pb.GetLotRequest) (entities.Lot, error)

	CreateAuction(ctx context.Context, in *pb.CreateAuctionRequest) (entities.Auction, error)
	CloseAuction(ctx context.Context, in *pb.CloseAuctionRequest) (entities.Auction, error)

	PlaceBid(ctx context.Context, in *pb.PlaceBidRequest) (entities.Bid, error)
	GetBid(ctx context.Context, in *pb.GetBidRequest) (entities.Bid, error)
}

func NewAuctionRepo(db database.DB, logger *slog.Logger) AuctionRepo {
	return &auctionRepo{
		db:     db,
		logger: logger,
	}
}
