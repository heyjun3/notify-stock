package notifystock

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func IsSameLen[T any](array ...[]T) bool {
	for i := range len(array) - 1 {
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
	currency := result[0].Meta.Currency

	stocks := make([]Stock, 0, len(timestamp))
	for i, t := range timestamp {
		stock, err := NewStock(
			symbol, time.Unix(int64(t), 0), currency,
			open[i], close[i], high[i], low[i],
		)
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

type MailService interface {
	Send(from, to, subject, text string) error
}

type StockNotifier struct {
	client      *FinanceClient
	mailService MailService
}

func NewStockNotifier(client *FinanceClient, mailService MailService) *StockNotifier {
	return &StockNotifier{
		client:      client,
		mailService: mailService,
	}
}

func (n *StockNotifier) Notify(symbols []string) error {
	now := time.Now()
	results := make([]Stocks, 0, len(symbols))
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
		results = append(results, stocks)
	}

	subject := fmt.Sprintf("Market Summary %s", now.Format("January 02 2006"))
	text := make([]string, 0)
	for _, result := range results {
		message, err := result.GenerateNotificationMessage()
		if err != nil {
			logger.Error("failed to generate notification message", "error", err)
			continue
		}
		text = append(text, message)
	}

	if err := n.mailService.Send(Cfg.FROM, Cfg.TO, subject, strings.Join(text, "\n\n")); err != nil {
		return err
	}
	return nil
}
