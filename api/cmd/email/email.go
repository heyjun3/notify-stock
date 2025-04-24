package email

import (
	"github.com/spf13/cobra"

	notifyStock "github.com/heyjun3/notify-stock/internal"
)

var EmailCommand = &cobra.Command{
	Use:   "email",
	Short: "send email.",
	Run: func(cmd *cobra.Command, args []string) {
		c := notifyStock.NewEmailClient(notifyStock.Cfg.MailToken)
		c.Send(notifyStock.Cfg.FROM, notifyStock.Cfg.TO, "test mail", "hello")
	},
}
