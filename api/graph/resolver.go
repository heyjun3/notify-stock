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
	symbolFetcher       *notify.SymbolFetcher
	notificationCreator *notify.NotificationCreator
	logger              *slog.Logger
}

func NewResolver(
	stockRepository *notify.StockRepository,
	symbolFetcher *notify.SymbolFetcher,
	notificationCreator *notify.NotificationCreator,
) *Resolver {
	return &Resolver{
		stockRepository:     stockRepository,
		symbolFetcher:       symbolFetcher,
		notificationCreator: notificationCreator,
		logger:              slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}
