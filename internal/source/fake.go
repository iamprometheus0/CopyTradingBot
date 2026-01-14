package source

import (
	"time"

	"github.com/iamprometheus0/CopyTradingBot/internal/types"
)

func StartFakeSource(out chan<- types.TradeEvent) {
	go func() {
		for {
			time.Sleep(2 * time.Second)

			out <- types.TradeEvent{
				MarketID: "BTC-15M-UP",
				Side:     "YES",
				Price:    0.52,
				Size:     100,
			}
		}
	}()
}
