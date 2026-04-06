package config

import (
	"os"
	"time"
)

type Config struct {
	Port            string
	UpstreamBaseURL string
	HTTPTimeout     time.Duration
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	upstream := os.Getenv("UPSTREAM_BASE_URL")
	if upstream == "" {
		upstream = "https://www.travel.taipei/open-api"
	}

	return Config{
		Port:            port,
		UpstreamBaseURL: upstream,
		HTTPTimeout:     10 * time.Second,
	}
}
