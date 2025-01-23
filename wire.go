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
		NewStockRegister,
	)
	return &StockRegister{}
}

func InitStockNotifier(ctx context.Context, credentialsPath string, client HTTPClientInterface) (*StockNotifier, error) {
	wire.Build(
		GmailServiceFactory,
		NewFinanceClient,
		NewStockNotifier,
	)
	return &StockNotifier{}, nil
}
