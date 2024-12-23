package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var config Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	config.FROM = os.Getenv("FROM")
	config.TO = os.Getenv("TO")
}

type Config struct {
	FROM string
	TO   string
}
