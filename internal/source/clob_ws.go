package source

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"github.com/iamprometheus0/CopyTradingBot/internal/config"
	"github.com/iamprometheus0/CopyTradingBot/internal/types"
)

type rtdsMsg struct {
	Event string `json:"event"`
	Data  struct {
		MarketID string  `json:"market_id"`
		Side     string  `json:"side"`
		Price    float64 `json:"price"`
		Size     float64 `json:"size"`
		Wallet   string  `json:"wallet"`
	} `json:"data"`
}

func StartCLOB(out chan<- types.TradeEvent, cfg *config.Config) {
	go func() {
		for {
			connectAndStream(out, cfg)
			log.Println("RTDS disconnected, reconnecting in 2s...")
			time.Sleep(2 * time.Second)
		}
	}()
}

func connectAndStream(out chan<- types.TradeEvent, cfg *config.Config) {
	c, _, err := websocket.DefaultDialer.Dial(
		"wss://ws-live-data.polymarket.com",
		nil,
	)
	if err != nil {
		log.Println("RTDS dial error:", err)
		return
	}
	defer c.Close()

	sub := map[string]any{
		"event":    "subscribe",
		"channels": []string{"trades"},
	}

	if err := c.WriteJSON(sub); err != nil {
		log.Println("RTDS subscribe error:", err)
		return
	}

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("RTDS read error:", err)
			return
		}

		var m rtdsMsg
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}

		if m.Event != "trade" {
			continue
		}

		//  ONLY COPY WHALE WALLETS
		if !cfg.WhaleWallets[m.Data.Wallet] {
			continue
		}

		log.Printf(
			"WHALE TRADE wallet=%s market=%s side=%s price=%.4f size=%.2f",
			m.Data.Wallet,
			m.Data.MarketID,
			m.Data.Side,
			m.Data.Price,
			m.Data.Size,
		)

		out <- types.TradeEvent{
			MarketID: m.Data.MarketID,
			Side:     m.Data.Side,
			Price:    m.Data.Price,
			Size:     m.Data.Size,
			Wallet:   m.Data.Wallet,
		}
	}
}
