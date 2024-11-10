package entities

type Lot struct {
	Id          int64
	Title       string
	Description string
	StartingBid float64
	SellerId    int64
	Status      string // Active, Closed
}
