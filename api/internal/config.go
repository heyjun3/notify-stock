package notifystock

import (
	"fmt"
	"log/slog"
	"os"

	yaml "github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
)

var Cfg Config

func init() {
	if err := initConfig(); err != nil {
		slog.Error("Failed to initialize configuration", "error", err)
		panic(err) // initではpanicが必要だが、エラー情報を改善
	}
}

func initConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../.env")
		if err != nil {
			slog.Info("No .env file found, using environment variables")
		}
	}

	cfg, err := loadConfigFromEnv()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	Cfg = *cfg
	logger = createLogger(Cfg.LogLevel)
	return nil
}

func loadConfigFromEnv() (*Config, error) {
	requiredEnvs := map[string]string{
		"FROM":                "",
		"TO":                  "",
		"MAIL_TOKEN":          "",
		"OAUTH_CLIENT_ID":     "",
		"OAUTH_CLIENT_SECRET": "",
		"OAUTH_REDIRECT_URL":  "",
	}

	// 必須環境変数の確認
	for key := range requiredEnvs {
		value, ok := os.LookupEnv(key)
		if !ok {
			return nil, fmt.Errorf("required environment variable %s is not set", key)
		}
		requiredEnvs[key] = value
	}

	// オプショナル環境変数
	dbdsn, ok := os.LookupEnv("DBDSN")
	if !ok {
		dbdsn = "postgres://postgres:postgres@localhost:5555/notify-stock?sslmode=disable"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	return &Config{
		FROM:              requiredEnvs["FROM"],
		TO:                requiredEnvs["TO"],
		DBDSN:             dbdsn,
		MailToken:         requiredEnvs["MAIL_TOKEN"],
		OauthClientID:     requiredEnvs["OAUTH_CLIENT_ID"],
		OauthClientSecret: requiredEnvs["OAUTH_CLIENT_SECRET"],
		OauthRedirectURL:  requiredEnvs["OAUTH_REDIRECT_URL"],
		LogLevel:          logLevel,
		Environment:       env,
	}, nil
}

type Config struct {
	FROM              string
	TO                string
	DBDSN             string
	MailToken         string
	OauthClientID     string
	OauthClientSecret string
	OauthRedirectURL  string
	LogLevel          string
	Environment       string
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
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
