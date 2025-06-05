package notifystock

import (
	"log/slog"
	"os"
	"strings"
)

var logger *slog.Logger

func CreateLogger(level string) *slog.Logger {
	var logLevel slog.Level
	l := strings.ToUpper(level)
	switch l {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "level" {
				a.Key = "severity"
			}
			if a.Key == "msg" {
				a.Key = "message"
			}
			return a
		},
	}))
}
