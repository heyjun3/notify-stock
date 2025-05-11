package notifystock

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

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

func (s *StockRegister) RegisterStockBySymbol(
	ctx context.Context, symbol string, start, end time.Time) error {
	stock, err := s.client.FetchStock(
		symbol,
		start,
		end,
		WithInterval("1d"),
	)
	if err != nil {
		return err
	}
	if err := s.stockRepository.Save(ctx, stock.stocks); err != nil {
		return err
	}
	if err := s.symbolRepository.Save(
		ctx, []SymbolDetail{stock.symbol},
	); err != nil {
		return err
	}
	return nil
}
func (s *StockRegister) RegisterStockBySymbols(
	ctx context.Context, symbols []string, start, end time.Time) error {
	var errs []error
	for _, symbol := range symbols {
		if err := s.RegisterStockBySymbol(
			ctx,
			symbol,
			start,
			end,
		); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type MailService interface {
	Send(from, to, subject, text string) error
}

type StockNotifier struct {
	mailService      MailService
	stockRepository  *StockRepository
	symbolRepository *SymbolRepository
}

func NewStockNotifier(
	mailService MailService,
	stockRepository *StockRepository,
	symbolRepository *SymbolRepository,
) *StockNotifier {
	return &StockNotifier{
		mailService:      mailService,
		stockRepository:  stockRepository,
		symbolRepository: symbolRepository,
	}
}

func (n *StockNotifier) Notify(symbols []string) error {
	ctx := context.Background()
	symbolDetails, err := n.symbolRepository.GetBySymbols(
		ctx, symbols,
	)
	if err != nil {
		return err
	}

	now := time.Now()
	stocks, err := n.stockRepository.GetStockByPeriodAndSymbols(
		ctx, symbols, now.AddDate(-1, 0, 0), now,
	)
	if err != nil {
		return err
	}
	results := make([]*Stocks, 0, len(symbols))
	for _, detail := range symbolDetails {
		stock, ok := stocks[detail.Symbol]
		if !ok {
			return fmt.Errorf("symbol %s not found", detail.Symbol)
		}
		result, err := NewStocks(detail, stock)
		if err != nil {
			return err
		}
		results = append(results, result)
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
