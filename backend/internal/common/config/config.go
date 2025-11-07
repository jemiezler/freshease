package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL               string
	HTTPPort                  string
	JWTSecret                 string
	Ent                       EntConfig
	OIDC_GOOGLE_ISSUER        string
	OIDC_GOOGLE_CLIENT_ID     string
	OIDC_GOOGLE_CLIENT_SECRET string
	OIDC_GOOGLE_REDIRECT_URI  string
	GENAI_APIKEY              string
	MinIO                     MinIOConfig
}

type EntConfig struct {
	Debug bool
}

type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	UseSSL          bool
}

// Load reads configuration from environment variables or defaults
func Load() Config {
	// Load .env file if it exists (useful for local dev)
	_ = godotenv.Load()

	cfg := Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:user1234@localhost:5432/trail-teller_db?sslmode=disable"),
		HTTPPort:    getEnv("HTTP_PORT", ":8080"),
		JWTSecret:   getEnv("JWT_SECRET", "secret"),
		Ent: EntConfig{
			Debug: getEnv("ENT_DEBUG", "false") == "true",
		},
		OIDC_GOOGLE_ISSUER:        getEnv("OIDC_GOOGLE_ISSUER", ""),
		OIDC_GOOGLE_CLIENT_ID:     getEnv("OIDC_GOOGLE_CLIENT_ID", ""),
		OIDC_GOOGLE_CLIENT_SECRET: getEnv("OIDC_GOOGLE_CLIENT_SECRET", ""),
		OIDC_GOOGLE_REDIRECT_URI:  getEnv("OIDC_GOOGLE_REDIRECT_URI", ""),
		GENAI_APIKEY:              getEnv("GENAI_APIKEY", ""),
		MinIO: MinIOConfig{
			Endpoint:        getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKeyID:     getEnv("MINIO_ACCESS_KEY_ID", "minioadmin"),
			SecretAccessKey: getEnv("MINIO_SECRET_ACCESS_KEY", "minioadmin1234"),
			Bucket:          getEnv("MINIO_BUCKET", "freshease"),
			UseSSL:          getEnv("MINIO_USE_SSL", "false") == "true",
		},
	}

	log.Printf("[config] Loaded config: DB=%s HTTP=%s EntDebug=%v", cfg.DatabaseURL, cfg.HTTPPort, cfg.Ent.Debug)
	return cfg
}

// getEnv returns the environment variable or a default if missing
func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
