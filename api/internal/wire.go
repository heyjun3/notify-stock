//go:build wireinject

package notifystock

import (
	"context"

	"github.com/google/wire"
)

func InitStockRegister(dsn string, client HTTPClientInterface) *StockRegister {
	wire.Build(
		NewDB,
		NewFinanceClient,
		NewStockRepository,
		NewSymbolRepository,
		NewStockRegister,
	)
	return &StockRegister{}
}

func InitSymbolFetcher(dsn string) *SymbolFetcher {
	wire.Build(
		NewDB,
		NewSymbolRepository,
		NewSymbolFetcher,
	)
	return &SymbolFetcher{}
}

func InitStockNotifier(ctx context.Context, token string, client HTTPClientInterface) (*StockNotifier, error) {
	wire.Build(
		NewEmailClient,
		NewFinanceClient,
		NewStockNotifier,
		wire.Bind(new(MailService), new(*EmailClient)),
	)
	return &StockNotifier{}, nil
}

func InitStockRepository(dsn string) *StockRepository {
	wire.Build(
		NewDB,
		NewStockRepository,
	)
	return &StockRepository{}
}

func InitNotificationRepository(dsn string) *NotificationRepository {
	wire.Build(
		NewDB,
		NewNotificationRepository,
	)
	return &NotificationRepository{}
}
