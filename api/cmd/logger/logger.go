package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

type MyHandler struct {
	slog.Handler
}

func NewMyHandler(h slog.Handler) *MyHandler {
	return &MyHandler{
		Handler: h,
	}
}
func (h *MyHandler) Handle(ctx context.Context, record slog.Record) error {
	// Custom handling logic can be added here
	record.AddAttrs(slog.Attr{Key: "custom", Value: slog.StringValue("custom_value")})
	return h.Handler.Handle(ctx, record)
}

var Command = &cobra.Command{
	Use:   "logger",
	Short: "Show logger use case",
	Run: func(cmd *cobra.Command, args []string) {
		log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
		log.Info("Logger command executed")
		log = slog.New(NewMyHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == "level" {
					a.Key = "severity"
				}
				if a.Key == "msg" {
					a.Key = "message"
				}
				return a
			},
		})))
		log.Info("Logger command executed with custom attributes")
	},
}
