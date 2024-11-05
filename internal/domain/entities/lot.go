package entities

type Lot struct {
	ID          int64
	Title       string
	Description string
	StartingBid float64
	SellerID    int64
	Status      string // Active, Closed
}
