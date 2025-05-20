//go:build wireinject

package graph

import (
	"github.com/google/wire"
	notify "github.com/heyjun3/notify-stock/internal"
	"log/slog"
)

func InitResolver(dsn string) *Resolver {
	wire.Build(
		notify.InitStockRepository,
		notify.InitNotificationRepository,
		notify.NewNotificationCreator,
		notify.NewNotificationFetcher,
		notify.InitSymbolFetcher,
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
