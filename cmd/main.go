package main

import (
	"github.com/heyjun3/notify-stock/cmd/notify"
	"github.com/heyjun3/notify-stock/cmd/stock"
	"github.com/heyjun3/notify-stock/cmd/version"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Short: "Stock notification Applications",
	}
)

func init() {
	rootCmd.AddCommand(version.VersionCommand, notify.NotifyCommand, stock.RegisterStockCommand)
}

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	Execute()
}
