package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	notify "github.com/heyjun3/notify-stock/internal"
	"github.com/spf13/cobra"
)

var FetchCommand = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch stock data",
	Run: func(cmd *cobra.Command, args []string) {
		if err := fetch(symbol); err != nil {
			fmt.Println("Error fetching stock data:", err)
			return
		}
		fmt.Println("Stock data fetched successfully")
	},
}

var symbol string

func init() {
	FetchCommand.Flags().StringVarP(&symbol, "symbol", "s", "", "Stock symbol to fetch data for")
}

func fetch(symbol string) error {
	s, err := notify.NewSymbol(symbol)
	if err != nil {
		return err
	}
	client := notify.NewFinanceClient(&http.Client{})
	res, err := client.FetchCurrentStock(s)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(res)
	if err != nil {
		return err
	}
	if err := os.WriteFile("stock_data.json", bytes, 0644); err != nil {
		return err
	}
	return nil
}
