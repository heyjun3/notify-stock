package register

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	notify "github.com/heyjun3/notify-stock/internal"
)

var logger *slog.Logger

var (
	symbols              []string
	isAll                bool
	RegisterStockCommand = &cobra.Command{
		Use:   "register",
		Short: "register stock by symbol and week",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if isAll {
				err = registerAllStockHistoryBySymbol(symbols)
			} else {
				err = registerStockByWeek(symbols)
			}
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	RegisterStockCommand.Flags().StringSliceVarP(&symbols, "symbol", "s", []string{}, "stock of symbol")
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
	return register.RegisterSymbols(symbols)
}

func registerAllStockHistoryBySymbol(symbols []string) error {
	register := notify.InitStockRegister(notify.Cfg.DBDSN, &http.Client{})
	end := time.Now().UTC()
	for {
		start := end.AddDate(-5, 0, 0)
		for _, symbol := range symbols {
			if err := register.SaveStock(symbol, start, end); err != nil {
				logger.Warn(err.Error())
				continue
			}
			time.Sleep(2 * time.Second)
		}
		end = start
		if end.Before(time.Unix(0, 0)) {
			break
		}
	}
	return register.RegisterSymbols(symbols)
}
