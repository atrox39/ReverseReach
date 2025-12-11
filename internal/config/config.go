package config

import (
	"os"

	"github.com/joho/godotenv"
)

type TunnelConfig struct {
	SSHUser       string
	SSHPassword   string
	SSHServerAddr string
	RemoteAddr    string
	LocalAddr     string
	WebPort       string
}

func Load() *TunnelConfig {
	godotenv.Load()

	cfg := &TunnelConfig{
		SSHUser:       getEnv("SSH_USER", "root"),
		SSHPassword:   getEnv("SSH_PASSWORD", "123456"),
		SSHServerAddr: getEnv("SSH_SERVER_ADDR", "example.com:22"),
		RemoteAddr:    getEnv("REMOTE_ADDR", "0.0.0.0:9001"),
		LocalAddr:     getEnv("LOCAL_ADDR", "127.0.0.1:9091"),
		WebPort:       getEnv("WEB_PORT", "3000"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
