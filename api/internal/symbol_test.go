package notifystock_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	notify "github.com/heyjun3/notify-stock/internal"
)

func TestSymbolDetail(t *testing.T) {
	t.Run("calculate change", func(t *testing.T) {
		detail := notify.NewSymbolDetail("N225", "N225", "Nikkei 225", "JPY",
			decimal.NewFromInt(1000), decimal.NewFromInt(900))

		assert.Equal(t, "+100", detail.Change())

		detail = notify.NewSymbolDetail("N225", "N225", "Nikkei 225", "JPY",
			decimal.NewFromInt(800), decimal.NewFromInt(900))

		assert.Equal(t, "-100", detail.Change())

		detail = notify.NewSymbolDetail("N225", "N225", "Nikkei 225", "JPY",
			decimal.NewFromInt(800), decimal.NewFromInt(800))

		assert.Equal(t, "0", detail.Change())
	})

	t.Run("calculate change percent", func(t *testing.T) {
		detail := notify.NewSymbolDetail("N225", "N225", "Nikkei 225", "JPY",
			decimal.NewFromInt(1000), decimal.NewFromInt(900))
		assert.Equal(t, "+11.11%", detail.ChangePercent())

		detail = notify.NewSymbolDetail("N225", "N225", "Nikkei 225", "JPY",
			decimal.NewFromInt(810), decimal.NewFromInt(900))

		assert.Equal(t, "-10%", detail.ChangePercent())

		detail = notify.NewSymbolDetail("N225", "N225", "Nikkei 225", "JPY",
			decimal.NewFromInt(900), decimal.NewFromInt(900))

		assert.Equal(t, "0%", detail.ChangePercent())
	})
}

func TestSymbolRepository(t *testing.T) {
	db := openDB(t)
	repo := notify.NewSymbolRepository(db)

	t.Run("save symbol", func(t *testing.T) {
		detail := notify.NewSymbolDetail("N225", "N225", "Nikkei 225", "JPY",
			decimal.NewFromInt(1000), decimal.NewFromInt(900))
		err := repo.Save(context.Background(), []notify.SymbolDetail{*detail})

		assert.NoError(t, err)
	})
	t.Run("update symbol", func(t *testing.T) {
		symbol := uuid.New().String()
		detail := notify.NewSymbolDetail(symbol, "N225", "Nikkei 225", "JPY",
			decimal.NewFromInt(1000), decimal.NewFromInt(900))
		err := repo.Save(context.Background(), []notify.SymbolDetail{*detail})
		assert.NoError(t, err)

		detail = notify.NewSymbolDetail(symbol, "N225"+symbol, "Nikkei 225"+symbol, "JPY",
			decimal.NewFromInt(1100), decimal.NewFromInt(1000))
		err = repo.Save(context.Background(), []notify.SymbolDetail{*detail})
		assert.NoError(t, err)

		s, err := repo.Get(context.Background(), symbol)

		assert.NoError(t, err)
		assert.Equal(t, symbol, s.Symbol)
		assert.Equal(t, "N225"+symbol, s.ShortName)
		assert.Equal(t, "Nikkei 225"+symbol, s.LongName)
		assert.Equal(t, decimal.NewFromInt(1100), s.MarketPrice)
		assert.Equal(t, decimal.NewFromInt(1000), s.PreviousClose)
	})

	t.Run("get symbol", func(t *testing.T) {
		detail := notify.NewSymbolDetail("S&P500", "S&P500", "S&P 500", "JPY",
			decimal.NewFromInt(1000), decimal.NewFromInt(900))
		err := repo.Save(context.Background(), []notify.SymbolDetail{*detail})
		assert.NoError(t, err)

		symbol, err := repo.Get(context.Background(), "S&P500")
		assert.NoError(t, err)
		assert.Equal(t, "S&P500", symbol.Symbol)
		assert.Equal(t, "S&P500", symbol.ShortName)
		assert.Equal(t, "S&P 500", symbol.LongName)
		assert.Equal(t, decimal.NewFromInt(1000), symbol.MarketPrice)
		assert.Equal(t, decimal.NewFromInt(900), symbol.PreviousClose)
		fmt.Println(symbol.Currency.String())
		assert.Equal(t, "JPY", symbol.Currency.String())
	})

	t.Run("get all symbols", func(t *testing.T) {
		detail := notify.NewSymbolDetail("N225", "N225", "Nikkei 225", "JPY",
			decimal.NewFromInt(1000), decimal.NewFromInt(900))
		err := repo.Save(context.Background(), []notify.SymbolDetail{*detail})
		assert.NoError(t, err)

		symbols, err := repo.GetAll(context.Background())
		assert.NoError(t, err)
		i := slices.IndexFunc(symbols, func(s notify.SymbolDetail) bool {
			return s.Symbol == "N225"
		})
		assert.Greater(t, i, -1)
	})
}
