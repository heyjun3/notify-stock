package main

import (
	"github.com/spf13/cobra"

	"github.com/heyjun3/notify-stock/cmd/email"
	"github.com/heyjun3/notify-stock/cmd/fetch"
	"github.com/heyjun3/notify-stock/cmd/notify"
	"github.com/heyjun3/notify-stock/cmd/register"
	"github.com/heyjun3/notify-stock/cmd/server"
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
		register.RegisterStockCommand,
		server.ServerCommand,
		email.EmailCommand,
		fetch.FetchCommand,
		yaml.YamlCommand,
	)
}

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	Execute()
}
