//go:build wireinject

package notifystock

import (
	"context"

	"github.com/google/wire"
	"github.com/uptrace/bun"
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

type DBDSN string

func wrapOpenDB(dsn DBDSN) *bun.DB {
	return NewDB(string(dsn))
}

func InitStockNotifier(
	ctx context.Context,
	token string,
	dsn DBDSN,
) (*StockNotifier, error) {
	wire.Build(
		wrapOpenDB,
		NewEmailClient,
		NewStockRepository,
		NewSymbolRepository,
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
