package entities

import "time"

type Bid struct {
	Id        int64
	AuctionId int64
	BidderId  int64
	Amount    float64
	CreatedAt time.Time
}
