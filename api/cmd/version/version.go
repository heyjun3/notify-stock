package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VersionCommand = &cobra.Command{
	Use:   "version",
	Short: "Show application version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version: 0.0.1")
	},
}
