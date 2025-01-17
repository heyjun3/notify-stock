package token

import (
	"log"

	"github.com/spf13/cobra"

	notify "github.com/heyjun3/notify-stock"
)

var (
	RefreshTokenCommand = &cobra.Command{
		Use:   "refresh-token",
		Short: "refresh notification token",
		Run: func(cmd *cobra.Command, args []string) {
			if err := notify.RefreshToken(); err != nil {
				log.Fatal(err)
			}
		},
	}
)
