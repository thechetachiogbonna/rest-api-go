package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type env struct {
	PORT               string
	DATABASE_URL       string
	JWT_SECRET         string
	JWT_REFRESH_SECRET string
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	return value
}

var Env = &env{
	PORT:               getEnv("PORT", "3000"),
	DATABASE_URL:       getEnv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/rest_api_go?sslmode=disable"),
	JWT_SECRET:         getEnv("JWT_SECRET", "f27786b1ddb611582156f7c3a4eb634a264c7b09e94ac6f9398b9e915709e2a5"),
	JWT_REFRESH_SECRET: getEnv("JWT_REFRESH_SECRET", "jwt_refresh_secret"),
}
