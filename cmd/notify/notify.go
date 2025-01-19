package notify

import (
	"net/http"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"

	notifyapp "github.com/heyjun3/notify-stock"
	"github.com/heyjun3/notify-stock/cmd/notify/token"
)

var NotifyCommand = &cobra.Command{
	Use:   "notify",
	Short: "Stock summary notification for yesterday",
	Run: func(cmd *cobra.Command, args []string) {
		notifyStock()
	},
}

func init() {
	NotifyCommand.AddCommand(token.RefreshTokenCommand)
}

type Stock struct {
	Symbol string
	Close  decimal.Decimal
	AVG    decimal.Decimal
	Start  time.Time
	End    time.Time
}

func GenerateStock(symbol notifyapp.Symbol, close []float64, timestamp []int) (*Stock, error) {
	s, err := symbol.ForFinance()
	if err != nil {
		return nil, err
	}
	avg, err := notifyapp.CalcAVG(close)
	if err != nil {
		return nil, err
	}
	latest := decimal.NewFromFloat(close[len(close)-1])
	start := time.Unix(int64(timestamp[0]), 0)
	end := time.Unix(int64(timestamp[len(timestamp)-1]), 0)
	return &Stock{
		Symbol: s,
		Close:  latest,
		AVG:    avg,
		Start:  start,
		End:    end,
	}, nil
}

func notifyStock() {
	notifier := notifyapp.InitStockNotifier(&http.Client{})
	notifier.Notify()
}
