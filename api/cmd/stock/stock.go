package stock

import (
	"github.com/spf13/cobra"

	"github.com/heyjun3/notify-stock/cmd/stock/update"
)

var (
	Command = &cobra.Command{
		Use:   "stock",
		Short: "Stock command",
	}
)

func init() {
	Command.AddCommand(
		update.StockCommand,
	)
}
