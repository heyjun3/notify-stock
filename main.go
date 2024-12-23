package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	client := NewFinanceClient(&http.Client{})
	now := time.Now()
	res, err := client.FetchStock("^N225", now.AddDate(0, -12, 0), now, WithInterval("1d"))
	if err != nil {
		log.Fatal(err)
	}
	close := res.Chart.Result[0].Indicators.Quote[0].Close
	avg, err := CalcAVG(close)
	if err != nil {
		log.Fatal(err)
	}
	subject := "Nikkei 225 Close & 1-Year MA"
	text := fmt.Sprintf("Nikkei 225 Data for %s\n", now.Format("January 02 2006")) +
		fmt.Sprintf("Closing Price: %d yen\n", int(close[len(close)-1])) +
		fmt.Sprintf("1-Year Moving Average: %v yen", avg.Ceil())
	err = NotifyGmail(context.Background(), config.FROM, config.TO, subject, text)
	if err != nil {
		slog.Error("error", "err", err)
	}
}
