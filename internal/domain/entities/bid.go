package entities

import "time"

type Bid struct {
	ID        int64
	AuctionID int64
	BidderID  int64
	Amount    float64
	CreatedAt time.Time
}
