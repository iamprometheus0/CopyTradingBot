package decision

import (
	"github.com/iamprometheus0/CopyTradingBot/internal/logging"
	"github.com/iamprometheus0/CopyTradingBot/internal/orderbook"
	"github.com/iamprometheus0/CopyTradingBot/internal/types"
)

const (
	maxAdversePct = 3.0
	maxSpreadPct  = 3.0
)

func Run(
	in <-chan types.TradeEvent,
	out chan<- types.Decision,
) {
	defer close(out)

	for evt := range in {
		book := orderbook.Get(evt.MarketID)

		spreadPct := SpreadMidPct(book.BestBid, book.BestAsk, book.Mid)
		if spreadPct > maxSpreadPct {
			logging.Logger.Printf(
				"SKIP instability: spread %.2f%%",
				spreadPct,
			)
			continue
		}

		prospective := book.BestAsk
		if evt.Side == "NO" {
			prospective = book.BestBid
		}

		dev := AdverseDeviationPct(evt.Price, prospective, evt.Side)
		if dev > maxAdversePct {
			logging.Logger.Printf(
				"SKIP slippage %.2f%%",
				dev,
			)
			continue
		}

		logging.Logger.Printf(
			"ACCEPT trade dev=%.2f%% spread=%.2f%%",
			dev,
			spreadPct,
		)

		out <- types.Decision{
			MarketID: evt.MarketID,
			Side:     evt.Side,
			Price:    prospective,
			Size:     evt.Size,
		}
	}
}
