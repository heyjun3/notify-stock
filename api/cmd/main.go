package main

import (
	"github.com/spf13/cobra"

	"github.com/heyjun3/notify-stock/cmd/email"
	"github.com/heyjun3/notify-stock/cmd/notify"
	"github.com/heyjun3/notify-stock/cmd/register"
	"github.com/heyjun3/notify-stock/cmd/server"
	"github.com/heyjun3/notify-stock/cmd/version"
)

var (
	rootCmd = &cobra.Command{
		Short: "Stock notification Applications",
	}
)

func init() {
	rootCmd.AddCommand(
		version.VersionCommand,
		notify.NotifyCommand,
		register.RegisterStockCommand,
		server.ServerCommand,
		email.EmailCommand,
	)
}

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	Execute()
}
