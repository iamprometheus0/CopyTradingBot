package source

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/iamprometheus0/CopyTradingBot/internal/types"
)

type clobMsg struct {
	Type string `json:"type"`
	Data struct {
		MarketID string  `json:"market_id"`
		Side     string  `json:"side"`
		Price    float64 `json:"price"`
		Size     float64 `json:"size"`
	} `json:"data"`
}

func StartCLOB(out chan<- types.TradeEvent) {
	c, _, err := websocket.DefaultDialer.Dial(
		"wss://ws-subscriptions-clob.polymarket.com/ws",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// subscribe to public trades
	sub := map[string]any{
		"type": "subscribe",
		"channels": []map[string]string{
			{"name": "trades"},
		},
	}
	_ = c.WriteJSON(sub)

	go func() {
		defer c.Close()
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("ws read error:", err)
				return
			}

			var m clobMsg
			if err := json.Unmarshal(msg, &m); err != nil {
				continue
			}

			if m.Type != "trade" {
				continue
			}

			log.Printf(
				"RAW TRADE market=%s side=%s price=%.4f size=%.2f",
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
			}
		}
	}()
}
