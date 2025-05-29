//go:build wireinject

package graph

import (
	"log/slog"

	"github.com/google/wire"
	"github.com/uptrace/bun"

	notify "github.com/heyjun3/notify-stock/internal"
)

func InitResolver(db *bun.DB) *Resolver {
	wire.Build(
		notify.InitStockRepository,
		notify.InitNotificationRepository,
		notify.InitSymbolRepository,
		notify.NewNotificationCreator,
		notify.NewNotificationFetcher,
		notify.NewDataLoader,
		NewResolver,
	)
	return &Resolver{}
}

func InitRootDirective(logger *slog.Logger) *DirectiveRoot {
	wire.Build(
		NewAuthDirective,
		NewDirectiveRoot,
	)
	return &DirectiveRoot{}
}
