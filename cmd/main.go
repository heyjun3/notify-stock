package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	notify "github.com/heyjun3/notify-stock"
	"github.com/shopspring/decimal"
)

const (
	N225  = "^N225"
	SP500 = "^GSPC"
)

type Stock struct {
	Symbol string
	Close  decimal.Decimal
	AVG    decimal.Decimal
	Start  time.Time
	End    time.Time
}

func GenerateStock(symbol string, close []float64, timestamp []int) (*Stock, error) {
	avg, err := notify.CalcAVG(close)
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

func main() {
	client := notify.NewFinanceClient(&http.Client{})
	now := time.Now()
	res, err := client.FetchStock(N225, now.AddDate(0, -12, 0), now, notify.WithInterval("1d"))
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
	res, err = client.FetchStock(SP500, now.AddDate(0, -12, 0), now, notify.WithInterval("1d"))
	if err != nil {
		log.Fatal(err)
	}
	sp500, err := GenerateStock(
		N225,
		res.Chart.Result[0].Indicators.Quote[0].Close,
		res.Chart.Result[0].Timestamp,
	)
	if err != nil {
		log.Fatal(err)
	}
	subject := fmt.Sprintf("Market Summary %s", now.Format("January 02 2006"))
	text := strings.Join([]string{
		"Nikkei 225 Data",
		fmt.Sprintf("Closing Price: %v yen", n225.Close.Ceil()),
		fmt.Sprintf("1-Year Moving Average: %v yen\n", n225.AVG.Ceil()),
		"S&P500 Data",
		fmt.Sprintf("Closing Price: %v $", sp500.Close.Ceil()),
		fmt.Sprintf("1-Year Moving Average: %v $", sp500.AVG.Ceil()),
	}, "\n")
	err = notify.NotifyGmail(context.Background(), notify.Cfg.FROM, notify.Cfg.TO, subject, text)
	if err != nil {
		slog.Error("error", "err", err)
	}
}
