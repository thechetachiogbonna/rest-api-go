package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type env struct {
	PORT         string
	DATABASE_URL string
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	return value
}

var Env = &env{
	PORT:         getEnv("PORT", "3000"),
	DATABASE_URL: getEnv("DATABASE_URL", ""),
}
