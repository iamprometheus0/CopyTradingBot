package config

import (
	"log"
	"os"
)

type Config struct {
	PrivateKey string
}

func Load() *Config {
	pk := os.Getenv("PRIVATE_KEY")
	if pk == "" {
		log.Fatal("PRIVATE_KEY not set")
	}

	return &Config{
		PrivateKey: pk,
	}
}
