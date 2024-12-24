package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var config Config

func init() {
	err := godotenv.Load()
	if err != nil {
		slog.Info(err.Error())
	}
	config.FROM = os.Getenv("FROM")
	config.TO = os.Getenv("TO")
}

type Config struct {
	FROM string
	TO   string
}
