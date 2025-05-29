//go:build wireinject

package notifystock

import (
	"context"

	"github.com/google/wire"
	"github.com/uptrace/bun"
)

func InitStockRegister(db *bun.DB, client HTTPClientInterface) *StockRegister {
	wire.Build(
		NewFinanceClient,
		NewStockRepository,
		NewSymbolRepository,
		NewStockRegister,
	)
	return &StockRegister{}
}

func InitSymbolRepository(db *bun.DB) *SymbolRepository {
	wire.Build(
		NewSymbolRepository,
	)
	return &SymbolRepository{}
}

func InitStockNotifier(
	ctx context.Context,
	token string,
	db *bun.DB,
) (*StockNotifier, error) {
	wire.Build(
		NewEmailClient,
		NewStockRepository,
		NewSymbolRepository,
		NewStockNotifier,
		wire.Bind(new(MailService), new(*EmailClient)),
	)
	return &StockNotifier{}, nil
}

func InitStockRepository(db *bun.DB) *StockRepository {
	wire.Build(
		NewStockRepository,
	)
	return &StockRepository{}
}

func InitNotificationRepository(db *bun.DB) *NotificationRepository {
	wire.Build(
		NewNotificationRepository,
	)
	return &NotificationRepository{}
}
