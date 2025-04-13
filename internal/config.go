package notifystock

import (
	"os"

	"github.com/joho/godotenv"
)

var Cfg Config

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		logger.Info(err.Error())
	}
	Cfg.FROM = os.Getenv("FROM")
	Cfg.TO = os.Getenv("TO")
	Cfg.DBDSN = "postgres://postgres:postgres@localhost:5555/notify-stock?sslmode=disable"
	Cfg.MailTrapToken = os.Getenv("MAIL_TRAP_TOKEN")
}

type Config struct {
	FROM          string
	TO            string
	DBDSN         string
	MailTrapToken string
}
