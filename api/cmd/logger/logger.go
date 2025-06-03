package logger

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "logger",
	Short: "Show logger use case",
	Run: func(cmd *cobra.Command, args []string) {
		log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
		log.Info("Logger command executed")
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
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
		log.Info("Logger command executed with custom attributes")
	},
}
