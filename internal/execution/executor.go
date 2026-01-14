package execution

import (
	"time"

	"github.com/iamprometheus0/CopyTradingBot/internal/logging"
	"github.com/iamprometheus0/CopyTradingBot/internal/types"
)

func Run(in <-chan types.Decision) {
	for d := range in {
		logging.Logger.Printf(
			"Executing order: %s %s @ %.2f size %.0f",
			d.MarketID,
			d.Side,
			d.Price,
			d.Size,
		)

		time.Sleep(300 * time.Millisecond)

		logging.Logger.Println("Order filled (simulated)")
	}
}
