package types

type TradeEvent struct {
	MarketID string
	Side     string
	Price    float64
	Size     float64
	Wallet   string
}
