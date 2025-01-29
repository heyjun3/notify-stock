package main

import (
	"github.com/heyjun3/notify-stock/cmd/notify"
	"github.com/heyjun3/notify-stock/cmd/register"
	"github.com/heyjun3/notify-stock/cmd/version"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Short: "Stock notification Applications",
	}
)

func init() {
	rootCmd.AddCommand(version.VersionCommand, notify.NotifyCommand, register.RegisterStockCommand)
}

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	Execute()
}
