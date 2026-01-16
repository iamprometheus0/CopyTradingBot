package config

import (
	"log"
	"os"
)

type Config struct {
	PrivateKey   string
	WhaleWallets map[string]bool
}

func Load() *Config {
	pk := os.Getenv("PRIVATE_KEY")
	if pk == "" {
		log.Fatal("PRIVATE_KEY not set")
	}

	return &Config{
		PrivateKey: pk,

		// 🔑 WHALE WALLETS (edit freely)
		WhaleWallets: map[string]bool{
			"0xabc123...": true,
			"0xdef456...": true,
		},
	}
}
