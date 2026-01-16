package source

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/iamprometheus0/CopyTradingBot/internal/types"
)

type rtdsMsg struct {
	Event string `json:"event"`
	Data  struct {
		MarketID string  `json:"market_id"`
		Side     string  `json:"side"`
		Price    float64 `json:"price"`
		Size     float64 `json:"size"`
	} `json:"data"`
}

func StartCLOB(out chan<- types.TradeEvent) {
	go func() {
		for {
			connectAndStream(out)
			log.Println("RTDS disconnected, reconnecting in 2s...")
			time.Sleep(2 * time.Second)
		}
	}()
}

func connectAndStream(out chan<- types.TradeEvent) {
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

		out <- types.TradeEvent{
			MarketID: m.Data.MarketID,
			Side:     m.Data.Side,
			Price:    m.Data.Price,
			Size:     m.Data.Size,
		}
	}
}
