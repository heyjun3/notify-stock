package notifystock

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func IsSameLen[T any](array ...[]T) bool {
	for i := 0; i < len(array)-1; i++ {
		if !(len(array[i]) == len(array[i+1])) {
			return false
		}
	}
	return true
}

func ConvertResponseToStock(res ChartResponse) ([]Stock, error) {
	result := res.Chart.Result
	if len(result) == 0 {
		return []Stock{}, fmt.Errorf("result is nil")
	}
	quote := result[0].Indicators.Quote
	if len(quote) == 0 {
		return []Stock{}, fmt.Errorf("quote is nil")
	}
	timestamp := result[0].Timestamp
	open := quote[0].Open
	close := quote[0].Close
	high := quote[0].High
	low := quote[0].Low
	if !IsSameLen(open, close, high, low) || len(timestamp) != len(open) {
		slog.Error(
			"same len error", "timestamp", len(timestamp), "open", len(open),
			"close", len(close), "high", len(high), "low", len(low))
		return []Stock{}, fmt.Errorf("don't same length error")
	}

	stocks := make([]Stock, len(timestamp))
	for i, t := range timestamp {
		stocks[i] = Stock{
			Timestamp: time.Unix(int64(t), 0),
			Open:      open[i],
			Close:     close[i],
			High:      high[i],
			Low:       low[i],
		}
	}
	return stocks, nil
}

func SaveStock(symbol string, begging, end time.Time) error {
	client := NewFinanceClient(&http.Client{})
	res, err := client.FetchStock(symbol, begging, end, WithInterval("1d"))
	if err != nil {
		return err
	}

}
