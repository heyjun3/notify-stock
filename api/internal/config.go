package notifystock

import (
	"os"

	yaml "github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
)

var Cfg Config

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../.env")
		if err != nil {
			logger.Info(err.Error())
		}
	}
	from, ok := os.LookupEnv("FROM")
	if !ok {
		panic("FROM is not set")
	}
	to, ok := os.LookupEnv("TO")
	if !ok {
		panic("TO is not set")
	}
	mailToken, ok := os.LookupEnv("MAIL_TOKEN")
	if !ok {
		panic("MAIL_TOKEN is not set")
	}
	dbdsn, ok := os.LookupEnv("DBDSN")
	if !ok {
		dbdsn = "postgres://postgres:postgres@localhost:5555/notify-stock?sslmode=disable"
	}
	Cfg = Config{
		FROM:      from,
		TO:        to,
		DBDSN:     dbdsn,
		MailToken: mailToken,
	}
}

type Config struct {
	FROM      string
	TO        string
	DBDSN     string
	MailToken string
}

type SupportSymbol struct {
	Symbols []string `yaml:"symbols"`
}

func GetSupportSymbols(path string) (*SupportSymbol, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var symbol SupportSymbol
	if err := yaml.Unmarshal(buf, &symbol); err != nil {
		return nil, err
	}
	return &symbol, nil
}
