package notifystock

import (
	"context"
	"fmt"
	"strings"
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

type StockNotifier struct {
	client *FinanceClient
}

func NewStockNotifier(client *FinanceClient) *StockNotifier {
	return &StockNotifier{
		client: client,
	}
}

func (n *StockNotifier) Notify() error {
	type StockWithSymbol struct {
		symbol Symbol
		stocks Stocks
	}
	symbols := []string{"N225", "S&P500"}
	now := time.Now()
	results := make([]StockWithSymbol, 0, len(symbols))
	for _, v := range symbols {
		symbol, err := NewSymbol(v)
		if err != nil {
			return err
		}
		res, err := n.client.FetchStock(symbol, now.AddDate(0, -12, 0), now, WithInterval("1d"))
		if err != nil {
			return err
		}
		stocks, err := ConvertResponseToStock(symbol, *res)
		if err != nil {
			return err
		}
		results = append(results, StockWithSymbol{symbol: symbol, stocks: stocks})
	}

	subject := fmt.Sprintf("Market Summary %s", now.Format("January 02 2006"))
	text := make([]string, 0, len(symbols)*3)
	for _, result := range results {
		avg, err := result.stocks.ClosingAverage()
		if err != nil {
			return err
		}
		latest := result.stocks.Latest()
		text = append(text,
			result.symbol.Display(),
			fmt.Sprintf("Closing Price: %v yen", int(latest.Close)),
			fmt.Sprintf("1-Year Moving Average: %v yen\n", avg.Ceil()),
		)
	}

	if err := NotifyGmail(context.Background(), Cfg.FROM, Cfg.TO, subject, strings.Join(text, "\n")); err != nil {
		return err
	}
	return nil
}
