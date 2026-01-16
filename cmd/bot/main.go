package main

import (
	"sync"

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

	tradeCh := make(chan types.TradeEvent, 1024)
	decisionCh := make(chan types.Decision, 1024)

	var wg sync.WaitGroup
	wg.Add(2)

	source.StartCLOB(tradeCh)

	go func() {
		defer wg.Done()
		decision.Run(tradeCh, decisionCh)
	}()

	go func() {
		defer wg.Done()
		execution.Run(decisionCh)
	}()

	wg.Wait()
}
