package utils

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string `env:"DB_URL,required"`
	Port      string `env:"PORT,required"`
	ENV       string `env:"ENV,required"`
}

// Cfg will hold your application’s config after Load()
var Cfg Config

// Load reads .env (if present) and then parses into Cfg.
func LoadENV() error {
	// ignore error if .env is missing (in prod you’ll use real ENV vars)
	_ = godotenv.Load()
	return env.Parse(&Cfg)
}
