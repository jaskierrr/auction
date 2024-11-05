package entities

type Auction struct {
	ID       int64
	LotID    int64
	Status   string // Active, Ended
	Bids     []Bid
	WinnerID *int64
}
