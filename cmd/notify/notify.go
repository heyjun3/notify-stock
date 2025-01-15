package notify

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"

	notifyapp "github.com/heyjun3/notify-stock"
)

var NotifyCommand = &cobra.Command{
	Use:   "notify",
	Short: "Stock summary notification for yesterday",
	Run: func(cmd *cobra.Command, args []string) {
		notifyStock()
	},
}

const (
	N225  = "^N225"
	SP500 = "^GSPC"
)

var SymbolMap = map[string]string{
	N225:  "Nikkei 225",
	SP500: "S&P 500",
}

type Stock struct {
	Symbol string
	Close  decimal.Decimal
	AVG    decimal.Decimal
	Start  time.Time
	End    time.Time
}

func GenerateStock(symbol string, close []float64, timestamp []int) (*Stock, error) {
	avg, err := notifyapp.CalcAVG(close)
	if err != nil {
		return nil, err
	}
	latest := decimal.NewFromFloat(close[len(close)-1])
	start := time.Unix(int64(timestamp[0]), 0)
	end := time.Unix(int64(timestamp[len(timestamp)-1]), 0)
	return &Stock{
		Symbol: symbol,
		Close:  latest,
		AVG:    avg,
		Start:  start,
		End:    end,
	}, nil
}

func notifyStock() {
	client := notifyapp.NewFinanceClient(&http.Client{})
	now := time.Now()
	res, err := client.FetchStock(N225, now.AddDate(0, -12, 0), now, notifyapp.WithInterval("1d"))
	if err != nil {
		log.Fatal(err)
	}
	n225, err := GenerateStock(
		N225,
		res.Chart.Result[0].Indicators.Quote[0].Close,
		res.Chart.Result[0].Timestamp,
	)
	if err != nil {
		log.Fatal(err)
	}
	res, err = client.FetchStock(SP500, now.AddDate(0, -12, 0), now, notifyapp.WithInterval("1d"))
	if err != nil {
		log.Fatal(err)
	}
	sp500, err := GenerateStock(
		SP500,
		res.Chart.Result[0].Indicators.Quote[0].Close,
		res.Chart.Result[0].Timestamp,
	)
	if err != nil {
		log.Fatal(err)
	}
	subject := fmt.Sprintf("Market Summary %s", now.Format("January 02 2006"))
	text := strings.Join([]string{
		SymbolMap[n225.Symbol],
		fmt.Sprintf("Closing Price: %v yen", n225.Close.Ceil()),
		fmt.Sprintf("1-Year Moving Average: %v yen\n", n225.AVG.Ceil()),
		SymbolMap[sp500.Symbol],
		fmt.Sprintf("Closing Price: %v $", sp500.Close.Ceil()),
		fmt.Sprintf("1-Year Moving Average: %v $", sp500.AVG.Ceil()),
	}, "\n")
	err = notifyapp.NotifyGmail(context.Background(), notifyapp.Cfg.FROM, notifyapp.Cfg.TO, subject, text)
	if err != nil {
		slog.Error("error", "err", err)
	}
}
