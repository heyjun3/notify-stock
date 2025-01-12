//go:build wireinject

package notifystock

import (
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
