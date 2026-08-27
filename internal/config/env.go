package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type env struct {
	PORT         string
	DATABASE_URL string
}

func getEnv() *env {
	return &env{
		PORT:         os.Getenv("PORT"),
		DATABASE_URL: os.Getenv("DATABASE_URL"),
	}
}

var Env = getEnv()
