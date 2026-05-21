package utils

import (
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string        `env:"PORT,required"`
	ENV                    string        `env:"ENV,required"`
	GRPCRequestTimeout     time.Duration `env:"GRPC_REQUEST_TIMEOUT" envDefault:"10s"`
	RoutingServiceAddr     string        `env:"ROUTING_SERVICE_ADDR,required"`
	AgentServiceAddr       string        `env:"AGENT_SERVICE_ADDR,required"`
	DbToolsAddr            string        `env:"DB_TOOLS_ADDR,required"`
	GeocodingAddr          string        `env:"GEOCODING_ADDR,required"`
	TrafficUpdaterAddr     string        `env:"TRAFFIC_UPDATER_ADDR,required"`
	DisableAuth            bool          `env:"DISABLE_AUTH" envDefault:"false"`
	SupabaseURL            string        `env:"SUPABASE_URL,required"`
	SupabaseSecretKey      string        `env:"SUPABASE_SECRET_KEY"`
	SupabaseServiceRoleKey string        `env:"SUPABASE_SERVICE_ROLE_KEY"`
}

// Cfg will hold your application’s config after Load()
var Cfg Config

// Load reads .env (if present) and then parses into Cfg.
func LoadENV() error {
	_ = godotenv.Load()
	return env.Parse(&Cfg)
}

// EffectiveSupabaseSecretKey returns the new secret key first, then the legacy service role key.
func (c Config) EffectiveSupabaseSecretKey() string {
	if v := strings.TrimSpace(c.SupabaseSecretKey); v != "" {
		return v
	}
	return strings.TrimSpace(c.SupabaseServiceRoleKey)
}
