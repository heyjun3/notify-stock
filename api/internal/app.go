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
		stocks, err := n.client.FetchStock(
			symbol,
			now.AddDate(0, -12, 0),
			now,
			WithInterval("1d"),
		)
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
