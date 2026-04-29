package utils

import (
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl              string `env:"DB_URL,required"`
	Port               string `env:"PORT,required"`
	ENV                string `env:"ENV,required"`
	RoutingServiceAddr string `env:"ROUTING_SERVICE_ADDR,required"`
	DbToolsAddr        string `env:"DB_TOOLS_ADDR,required"`
	GeocodingAddr      string `env:"GEOCODING_ADDR,required"`
	TrafficUpdaterAddr string `env:"TRAFFIC_UPDATER_ADDR,required"`
	DisableAuth        bool   `env:"DISABLE_AUTH" envDefault:"false"`
	FirebaseProjectID  string `env:"FIREBASE_PROJECT_ID"`
	GoogleCloudProject string `env:"GOOGLE_CLOUD_PROJECT"`
}

// Cfg will hold your application’s config after Load()
var Cfg Config

// Load reads .env (if present) and then parses into Cfg.
func LoadENV() error {
	_ = godotenv.Load()
	return env.Parse(&Cfg)
}

// EffectiveFirebaseProjectID returns FIREBASE_PROJECT_ID first, then GOOGLE_CLOUD_PROJECT.
func (c Config) EffectiveFirebaseProjectID() string {
	if v := strings.TrimSpace(c.FirebaseProjectID); v != "" {
		return v
	}
	return strings.TrimSpace(c.GoogleCloudProject)
}
