package notifystock

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func IsSameLen[T any](array ...[]T) bool {
	for i := range len(array) - 1 {
		if !(len(array[i]) == len(array[i+1])) {
			return false
		}
	}
	return true
}

func ConvertResponseToStock(res ChartResponse) (*Stocks, error) {
	result := res.Chart.Result
	if len(result) == 0 {
		return nil, fmt.Errorf("result is nil")
	}
	quote := result[0].Indicators.Quote
	if len(quote) == 0 {
		return nil, fmt.Errorf("quote is nil")
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
		return nil, fmt.Errorf("don't same length error")
	}
	currency := result[0].Meta.Currency
	symbol := result[0].Meta.Symbol

	stocks := make([]Stock, 0, len(timestamp))
	for i, t := range timestamp {
		stock, err := NewStock(
			symbol, time.Unix(int64(t), 0),
			open[i], close[i], high[i], low[i],
		)
		if err != nil {
			logger.Error("new stock error", "error", err)
			continue
		}
		stocks = append(stocks, stock)
	}
	detail, err := ConvertResponseToSymbol(&res)
	if err != nil {
		return nil, err
	}
	s, err := NewStocks(*detail, currency, stocks)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func ConvertResponseToSymbol(res *ChartResponse) (*SymbolDetail, error) {
	result := res.Chart.Result
	if len(result) == 0 {
		return nil, fmt.Errorf("result is nil")
	}
	meta := result[0].Meta
	detail := NewSymbolDetail(meta.Symbol, meta.ShortName, meta.LongName, meta.Currency,
		decimal.NewFromFloat(meta.RegularMarketPrice), decimal.NewFromFloat(meta.ChartPreviousClose))
	return detail, nil
}

type StockRegister struct {
	client           *FinanceClient
	stockRepository  *StockRepository
	symbolRepository *SymbolRepository
}

func NewStockRegister(
	client *FinanceClient,
	stockRepository *StockRepository,
	symbolRepository *SymbolRepository,
) *StockRegister {
	return &StockRegister{
		client:           client,
		stockRepository:  stockRepository,
		symbolRepository: symbolRepository,
	}
}

func (s *StockRegister) SaveStock(symbol string, begging, end time.Time) error {
	ctx := context.Background()
	res, err := s.client.FetchStock(symbol, begging, end, WithInterval("1d"))
	if err != nil {
		return err
	}
	stock, err := ConvertResponseToStock(*res)
	if err != nil {
		return err
	}
	if err := s.stockRepository.Save(ctx, stock.stocks); err != nil {
		return err
	}
	return nil
}

func (s *StockRegister) RegisterSymbol(symbol string) error {
	ctx := context.Background()
	res, err := s.client.FetchCurrentStock(symbol)
	if err != nil {
		return err
	}
	detail, err := ConvertResponseToSymbol(res)
	if err != nil {
		return err
	}
	if err := s.symbolRepository.Save(ctx, []SymbolDetail{*detail}); err != nil {
		return err
	}
	return nil
}

func (s *StockRegister) RegisterSymbols(symbols []string) error {
	for _, symbol := range symbols {
		if err := s.RegisterSymbol(symbol); err != nil {
			return err
		}
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
	results := make([]*Stocks, 0, len(symbols))
	for _, symbol := range symbols {
		res, err := n.client.FetchStock(symbol, now.AddDate(0, -12, 0), now, WithInterval("1d"))
		if err != nil {
			return err
		}
		stocks, err := ConvertResponseToStock(*res)
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
