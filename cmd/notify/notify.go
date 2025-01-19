package notify

import (
	// "context"
	// "fmt"
	// "log"
	"net/http"
	// "strings"
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

	// client := notifyapp.NewFinanceClient(&http.Client{})
	// n225Symbol, err := notifyapp.NewSymbol("N225")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// now := time.Now()
	// res, err := client.FetchStock(n225Symbol, now.AddDate(0, -12, 0), now, notifyapp.WithInterval("1d"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// n225, err := GenerateStock(
	// 	n225Symbol,
	// 	res.Chart.Result[0].Indicators.Quote[0].Close,
	// 	res.Chart.Result[0].Timestamp,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// spSymbol, err := notifyapp.NewSymbol("S&P500")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// res, err = client.FetchStock(spSymbol, now.AddDate(0, -12, 0), now, notifyapp.WithInterval("1d"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// sp500, err := GenerateStock(
	// 	spSymbol,
	// 	res.Chart.Result[0].Indicators.Quote[0].Close,
	// 	res.Chart.Result[0].Timestamp,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// subject := fmt.Sprintf("Market Summary %s", now.Format("January 02 2006"))
	// text := strings.Join([]string{
	// 	n225Symbol.Display(),
	// 	fmt.Sprintf("Closing Price: %v yen", n225.Close.Ceil()),
	// 	fmt.Sprintf("1-Year Moving Average: %v yen\n", n225.AVG.Ceil()),
	// 	spSymbol.Display(),
	// 	fmt.Sprintf("Closing Price: %v $", sp500.Close.Ceil()),
	// 	fmt.Sprintf("1-Year Moving Average: %v $", sp500.AVG.Ceil()),
	// }, "\n")
	// err = notifyapp.NotifyGmail(context.Background(), notifyapp.Cfg.FROM, notifyapp.Cfg.TO, subject, text)
	// if err != nil {
	// 	log.Fatal("error", "err", err)
	// }
}
