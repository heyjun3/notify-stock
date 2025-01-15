package cmd

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
	rootCmd.AddCommand(version.VersionCommand)
	rootCmd.AddCommand(notify.NotifyCommand)
	rootCmd.AddCommand(stock.RegisterStockByWeekCommand)
}

func Execute() error {
	return rootCmd.Execute()
}
