package notifystock_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	notify "github.com/heyjun3/notify-stock/internal"
)

func TestDataLoaderSymbolDetail(t *testing.T) {
	db := openDB(t)
	symbolRepository := notify.NewSymbolRepository(db)
	dataloader := notify.NewDataLoader(symbolRepository)

	t.Run("Load Existing Key", func(t *testing.T) {
		symbol := notify.NewSymbolDetail(
			"N225",
			"日経225",
			"日経平均株価",
			"JPY",
			decimal.NewFromInt(10000),
			decimal.NewFromInt(9900),
		)
		err := symbolRepository.Save(context.Background(), []notify.SymbolDetail{*symbol})
		assert.NoError(t, err)

		result, err := dataloader.SymbolDetail.Load(context.Background(), "N225")()

		assert.NoError(t, err)
		assert.Equal(t, symbol.Symbol, result.Symbol)
		assert.Equal(t, symbol.ShortName, result.ShortName)
		assert.Equal(t, symbol.LongName, result.LongName)
		assert.Equal(t, symbol.Currency.Symbol(), result.Currency.Symbol())
	})

	t.Run("Not Found Key", func(t *testing.T) {
		result, err := dataloader.SymbolDetail.Load(context.Background(), uuid.New().String())()

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
