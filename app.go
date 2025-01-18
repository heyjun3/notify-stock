package notifystock

import (
	"context"
	"fmt"
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

func ConvertResponseToStock(symbol Symbol, res ChartResponse) ([]Stock, error) {
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
		logger.Error(
			"same len error", "timestamp", len(timestamp), "open", len(open),
			"close", len(close), "high", len(high), "low", len(low))
		return []Stock{}, fmt.Errorf("don't same length error")
	}

	stocks := make([]Stock, 0, len(timestamp))
	for i, t := range timestamp {
		stock, err := NewStock(symbol, time.Unix(int64(t), 0), open[i], close[i], high[i], low[i])
		if err != nil {
			logger.Error("new stock error", "error", err)
			continue
		}
		stocks = append(stocks, stock)
	}
	return stocks, nil
}

type StockRegister struct {
	client          *FinanceClient
	stockRepository *StockRepository
}

func NewStockRegister(client *FinanceClient, repository *StockRepository) *StockRegister {
	return &StockRegister{
		client:          client,
		stockRepository: repository,
	}
}

func (s *StockRegister) SaveStock(symbol string, begging, end time.Time) error {
	sym, err := NewSymbol(symbol)
	if err != nil {
		return err
	}
	res, err := s.client.FetchStock(sym, begging, end, WithInterval("1d"))
	if err != nil {
		return err
	}
	stocks, err := ConvertResponseToStock(sym, *res)
	if err != nil {
		return err
	}
	if err := s.stockRepository.Save(context.Background(), stocks); err != nil {
		return err
	}
	return nil
}
