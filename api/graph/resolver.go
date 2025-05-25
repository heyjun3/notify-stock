package graph

import (
	"log/slog"
	"os"

	notify "github.com/heyjun3/notify-stock/internal"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	stockRepository     *notify.StockRepository
	symbolRepository    *notify.SymbolRepository
	notificationCreator *notify.NotificationCreator
	notificationFetcher *notify.NotificationFetcher
	logger              *slog.Logger
	loader              *notify.DataLoader
}

func NewResolver(
	stockRepository *notify.StockRepository,
	symbolRepository *notify.SymbolRepository,
	notificationCreator *notify.NotificationCreator,
	notificationFetcher *notify.NotificationFetcher,
	loader *notify.DataLoader,
) *Resolver {
	return &Resolver{
		stockRepository:     stockRepository,
		symbolRepository:    symbolRepository,
		notificationCreator: notificationCreator,
		notificationFetcher: notificationFetcher,
		logger:              slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		loader:              loader,
	}
}
