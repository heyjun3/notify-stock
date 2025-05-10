package fetch

import (
	"fmt"
	"net/http"
	"os"
	"time"

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
	client := notify.NewFinanceClient(&http.Client{})
	// res, err := client.FetchCurrentStock(symbol)
	// res, err := client.FetchStock(symbol, time.Now().AddDate(0, 0, -7), time.Now(), notify.WithInterval("1d"))
	res, err := client.FetchStock(symbol, time.Now().AddDate(0, 0, -7), time.Now())
	if err != nil {
		return err
	}
	bytes, err := res.Json()
	if err != nil {
		return err
	}
	if err := os.WriteFile("stock_data.json", bytes, 0644); err != nil {
		return err
	}
	return nil
}
