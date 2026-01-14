package decision

import (
	"github.com/iamprometheus0/CopyTradingBot/internal/logging"
	"github.com/iamprometheus0/CopyTradingBot/internal/types"
)

func Run(
	in <-chan types.TradeEvent,
	out chan<- types.Decision,
) {
	for evt := range in {
		logging.Logger.Printf(
			"Decision: copying %s %s @ %.2f size %.0f",
			evt.MarketID,
			evt.Side,
			evt.Price,
			evt.Size,
		)

		out <- types.Decision{
			MarketID: evt.MarketID,
			Side:     evt.Side,
			Price:    evt.Price,
			Size:     evt.Size,
		}
	}
}
