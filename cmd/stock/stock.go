package stock

import (
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	notify "github.com/heyjun3/notify-stock"
)

var (
	symbol               string
	isAll                bool
	RegisterStockCommand = &cobra.Command{
		Use:   "register",
		Short: "register stock by symbol and week",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if isAll {
				err = registerAllStockHistoryBySymbol(symbol)
			} else {
				err = registerStockByWeek(symbol)
			}
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	RegisterStockCommand.Flags().StringVarP(&symbol, "symbol", "s", "", "stock of symbol")
	RegisterStockCommand.Flags().BoolVarP(&isAll, "all", "a", false, "register stock price data for the entire period")
	RegisterStockCommand.MarkFlagRequired("symbol")
}

func registerStockByWeek(symbol string) error {
	register := notify.InitStockRegister(notify.Cfg.DBDSN, &http.Client{})
	now := time.Now().UTC()
	if err := register.SaveStock(symbol, now.AddDate(0, 0, -7), now); err != nil {
		return err
	}
	return nil
}

func registerAllStockHistoryBySymbol(symbol string) error {
	register := notify.InitStockRegister(notify.Cfg.DBDSN, &http.Client{})
	t := time.Unix(0, 0)
	times := []time.Time{t}
	for {
		t = t.AddDate(5, 0, 0)
		if t.After(time.Now().UTC()) {
			t = time.Now().UTC()
			times = append(times, t)
			break
		}
		times = append(times, t)
	}
	for i := 0; i < len(times)-1; i++ {
		if err := register.SaveStock(
			symbol, times[i], times[i+1]); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}
