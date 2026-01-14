package main

import (
	"github.com/joho/godotenv"

	"github.com/iamprometheus0/CopyTradingBot/internal/config"
	"github.com/iamprometheus0/CopyTradingBot/internal/decision"
	"github.com/iamprometheus0/CopyTradingBot/internal/execution"
	"github.com/iamprometheus0/CopyTradingBot/internal/logging"
	"github.com/iamprometheus0/CopyTradingBot/internal/source"
	"github.com/iamprometheus0/CopyTradingBot/internal/types"
)

func main() {
	logging.Init()

	godotenv.Load()
	config.Load()

	tradeCh := make(chan types.TradeEvent, 10)
	decisionCh := make(chan types.Decision, 10)

	source.StartFakeSource(tradeCh)
	go decision.Run(tradeCh, decisionCh)
	execution.Run(decisionCh)
}
