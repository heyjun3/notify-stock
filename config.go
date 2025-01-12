package notifystock

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var Cfg Config

func init() {
	err := godotenv.Load()
	if err != nil {
		slog.Info(err.Error())
	}
	Cfg.FROM = os.Getenv("FROM")
	Cfg.TO = os.Getenv("TO")
	Cfg.DBDSN = "postgres://postgres:postgres@localhost:5555/notify-stock?sslmode=disable"
}

type Config struct {
	FROM  string
	TO    string
	DBDSN string
}
