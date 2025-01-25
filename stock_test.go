package notifystock_test

import (
	"context"
	"testing"
	"time"

	notify "github.com/heyjun3/notify-stock"
	"github.com/stretchr/testify/assert"

	"github.com/uptrace/bun"
)

var db *bun.DB

func init() {
	dsn := "postgres://postgres:postgres@localhost:5555/notify-stock-test?sslmode=disable"
	db = notify.NewDB(dsn)
}

func TestSave(t *testing.T) {
	repo := notify.NewStockRepository(db)

	tests := []struct {
		name   string
		stocks []notify.Stock
		err    error
	}{
		{
			stocks: []notify.Stock{{Symbol: "N255", Timestamp: time.Now(), Open: 1000, Close: 2000, High: 2500, Low: 500}},
			err:    nil,
		},
		{
			stocks: []notify.Stock{},
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Save(context.Background(), tt.stocks)

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestGetStockByPeriod(t *testing.T) {
	repo := notify.NewStockRepository(db)

	tests := []struct {
		name   string
		symbol func() notify.Symbol
		err    error
	}{{
		symbol: func() notify.Symbol {
			symbol, _ := notify.NewSymbol("N225")
			return symbol
		},
		err: nil,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.GetStockByPeriod(context.Background(), tt.symbol(), time.Now(), time.Now())

			assert.NoError(t, err)
		})
	}
}

func TestStocksLatest(t *testing.T) {
	stocks := notify.Stocks{
		{
			Timestamp: time.Now(),
			Close:     10000,
		},
		{
			Timestamp: time.Now().AddDate(-1, 0, 0),
			Close:     90000,
		},
	}

	t.Run("", func(t *testing.T) {
		latest := stocks.Latest()

		assert.Equal(t, float64(10000), latest.Close)
	})
}
