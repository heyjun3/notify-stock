package notify

import (
	"context"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	notifyapp "github.com/heyjun3/notify-stock"
	"github.com/heyjun3/notify-stock/cmd/notify/token"
)

var NotifyCommand = &cobra.Command{
	Use:   "notify",
	Short: "Stock summary notification for yesterday",
	Run: func(cmd *cobra.Command, args []string) {
		notifyStock(symbols)
	},
}

var symbols []string

func init() {
	NotifyCommand.AddCommand(token.RefreshTokenCommand)
	NotifyCommand.Flags().StringSliceVarP(&symbols, "symbol", "s", []string{}, "array of symbol")
}

func notifyStock(symbols []string) {
	notifier, err := notifyapp.InitStockNotifier(context.Background(), notifyapp.Cfg.MailTrapToken, &http.Client{})
	if err != nil {
		log.Fatal(err)
	}
	if err := notifier.Notify(symbols); err != nil {
		log.Fatal(err)
	}
}
