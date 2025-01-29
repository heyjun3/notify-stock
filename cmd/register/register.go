package register

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	notify "github.com/heyjun3/notify-stock"
)

var logger *slog.Logger

var (
	symbol               []string
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
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	RegisterStockCommand.Flags().StringSliceVarP(&symbol, "symbol", "s", []string{}, "stock of symbol")
	RegisterStockCommand.Flags().BoolVarP(&isAll, "all", "a", false, "register stock price data for the entire period")
	RegisterStockCommand.MarkFlagRequired("symbol")
}

func registerStockByWeek(symbols []string) error {
	register := notify.InitStockRegister(notify.Cfg.DBDSN, &http.Client{})
	now := time.Now().UTC()
	for _, symbol := range symbols {
		if err := register.SaveStock(symbol, now.AddDate(0, 0, -7), now); err != nil {
			logger.Error(err.Error())
		}
	}
	return nil
}

func registerAllStockHistoryBySymbol(symbols []string) error {
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
		for _, symbol := range symbols {
			if err := register.SaveStock(
				symbol, times[i], times[i+1]); err != nil {
				return err
			}
			time.Sleep(2 * time.Second)
		}
	}
	return nil
}
