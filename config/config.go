package config

import (
	"time"
)

// Config holds application configuration
type Config struct {
	Server ServerConfig
	Gacha  GachaConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port            string
	AllowedOrigins  []string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// GachaConfig holds gacha system configuration
type GachaConfig struct {
	SinglePullCost int
	TenPullCost    int
	PityThreshold  int
	SSRRate        float64
	SRRate         float64
	RRate          float64
}

// LoadConfig loads configuration (can be extended to read from file/env)
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            ":8080",
			AllowedOrigins:  []string{"*"},
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			ShutdownTimeout: 10 * time.Second,
		},
		Gacha: GachaConfig{
			SinglePullCost: 160,
			TenPullCost:    1600,
			PityThreshold:  90,
			SSRRate:        0.02, // 2%
			SRRate:         0.10, // 10%
			RRate:          0.88, // 88%
		},
	}
}
