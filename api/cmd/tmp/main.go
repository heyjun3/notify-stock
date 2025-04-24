package main

import (
	"log/slog"
	"os"
	"time"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func main() {
	logger.Info("test")
	t := time.Now()
	logger.Info(t.Format(time.RFC3339))

	tt := t.Round(time.Hour).UTC()
	logger.Info(tt.Format(time.RFC3339))
}
