package entities

type Auction struct {
	Id      int64
	LotId    int64
	Status   string // Active, Ended
	Bids     []Bid
	WinnerId *int64
}
