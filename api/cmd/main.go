package main

import (
	"github.com/spf13/cobra"

	"github.com/heyjun3/notify-stock/cmd/email"
	"github.com/heyjun3/notify-stock/cmd/fetch"
	"github.com/heyjun3/notify-stock/cmd/logger"
	"github.com/heyjun3/notify-stock/cmd/notify"
	"github.com/heyjun3/notify-stock/cmd/server"
	"github.com/heyjun3/notify-stock/cmd/stock"
	"github.com/heyjun3/notify-stock/cmd/version"
	"github.com/heyjun3/notify-stock/cmd/yaml"
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
		server.ServerCommand,
		email.EmailCommand,
		fetch.FetchCommand,
		yaml.YamlCommand,
		stock.Command,
		logger.Command,
	)
}

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	Execute()
}
