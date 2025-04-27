//go:build wireinject

package graph

import (
	"github.com/google/wire"
	notify "github.com/heyjun3/notify-stock/internal"
)

func InitResolver(dsn string) *Resolver {
	wire.Build(
		notify.InitStockRepository,
		notify.InitNotificationRepository,
		notify.NewNotificationCreator,
		notify.InitSymbolFetcher,
		NewResolver,
	)
	return &Resolver{}
}
