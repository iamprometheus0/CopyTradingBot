package types

type TradeEvent struct {
	MarketID string
	Side     string // "YES" or "NO"
	Price    float64
	Size     float64
}
