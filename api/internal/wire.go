//go:build wireinject

package notifystock

import (
	"context"
	"net/http"

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
	db *bun.DB,
	config MailGunClientConfig,
) (*StockNotifier, error) {
	wire.Build(
		NewMailGunClient,
		NewStockRepository,
		NewSymbolRepository,
		NewStockNotifier,
		wire.Bind(new(MailService), new(*MailGunClient)),
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

func InitAuthHandler(sessions *Sessions, db *bun.DB, client http.Client, option GoogleClientOption) *AuthHandler {
	wire.Build(
		NewGoogleClient,
		NewMemberRepository,
		NewAuthHandler,
	)
	return &AuthHandler{}
}

func InitNotificationCreator(db *bun.DB) *NotificationCreator {
	wire.Build(
		NewNotificationRepository,
		NewSymbolRepository,
		NewNotificationCreator,
	)
	return &NotificationCreator{}
}
