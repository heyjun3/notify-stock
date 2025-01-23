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
		notifyStock()
	},
}

func init() {
	NotifyCommand.AddCommand(token.RefreshTokenCommand)
}

func notifyStock() {
	notifier, err := notifyapp.InitStockNotifier(context.Background(), "credentials.json", &http.Client{})
	if err != nil {
		log.Fatal(err)
	}
	if err := notifier.Notify([]string{}); err != nil {
		log.Fatal(err)
	}
}
