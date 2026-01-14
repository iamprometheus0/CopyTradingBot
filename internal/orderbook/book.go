package orderbook

import "math/rand"

type Book struct {
	BestBid float64
	BestAsk float64
	Mid     float64
}

func Get(marketID string) Book {
	mid := 0.52 + (rand.Float64()-0.5)*0.02
	spread := 0.01

	return Book{
		BestBid: mid - spread/2,
		BestAsk: mid + spread/2,
		Mid:     mid,
	}
}
