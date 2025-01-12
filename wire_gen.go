// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package notifystock

// Injectors from wire.go:

func InitStockRegister(dsn string, client HTTPClientInterface) *StockRegister {
	financeClient := NewFinanceClient(client)
	db := NewDB(dsn)
	stockRepository := NewStockRepository(db)
	stockRegister := NewStockRegister(financeClient, stockRepository)
	return stockRegister
}
