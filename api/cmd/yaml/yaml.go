package yaml

import (
	"fmt"
	"os"

	y "github.com/goccy/go-yaml"
	"github.com/spf13/cobra"

	notify "github.com/heyjun3/notify-stock/internal"
)

var YamlCommand = &cobra.Command{
	Use:   "yaml",
	Short: "parse yaml file",
	Run: func(cmd *cobra.Command, args []string) {
		buf, err := os.ReadFile("config.yaml")
		if err != nil {
			panic(err)
		}
		var symbol notify.SupportSymbol
		if err := y.Unmarshal(buf, &symbol); err != nil {
			panic(err)
		}
		for _, symbol := range symbol.Symbols {
			fmt.Println("symbol", symbol)
		}
	},
}
