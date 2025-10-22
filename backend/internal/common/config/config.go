package config

import "os"

type Config struct {
	HTTPAddr  string
	JWTSecret string
}

func Load() Config {
	cfg := Config{
		HTTPAddr:  env("HTTP_ADDR", ":8080"),
		JWTSecret: env("JWT_SECRET", "secret"),
	}
	return cfg
}

func env(k, def string) string {
	if val := os.Getenv(k); val != "" {
		return val
	}
	return def
}
