package email

import (
	"github.com/spf13/cobra"

	notifyStock "github.com/heyjun3/notify-stock"
)

var EmailCommand = &cobra.Command{
	Use:   "email",
	Short: "send email.",
	Run: func(cmd *cobra.Command, args []string) {
		c := notifyStock.NewMailTrapClient(notifyStock.Cfg.MailTrapToken)
		c.Send("hello@demomailtrap.co", notifyStock.Cfg.TO, "test mail", "hello")
	},
}
