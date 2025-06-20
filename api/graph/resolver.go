package graph

import (
	"log/slog"

	notify "github.com/heyjun3/notify-stock/internal"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	stockRepository        *notify.StockRepository
	symbolRepository       *notify.SymbolRepository
	notificationRepository *notify.NotificationRepository
	logger                 *slog.Logger
	loader                 *notify.DataLoader
}

func NewResolver(
	stockRepository *notify.StockRepository,
	symbolRepository *notify.SymbolRepository,
	notificationRepository *notify.NotificationRepository,
	loader *notify.DataLoader,
) *Resolver {
	return &Resolver{
		stockRepository:        stockRepository,
		symbolRepository:       symbolRepository,
		notificationRepository: notificationRepository,
		logger:                 notify.CreateLogger("info"),
		loader:                 loader,
	}
}
