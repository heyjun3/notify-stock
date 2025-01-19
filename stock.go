package notifystock

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewDB(dsn string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))
	return db
}

func NewStock(symbol Symbol, timestamp time.Time,
	open, close, high, low float64) (Stock, error) {
	s, err := symbol.ForDB()
	if err != nil {
		return Stock{}, err
	}
	for _, v := range []float64{open, close, high, low} {
		if v <= 0 {
			return Stock{}, fmt.Errorf(
				"value is higher than zero. open: %v, close: %v, high: %v, low: %v", open, close, high, low)
		}
	}
	return Stock{
		Symbol:    s,
		Timestamp: timestamp,
		Open:      open,
		Close:     close,
		High:      high,
		Low:       low,
	}, nil
}

type Stock struct {
	bun.BaseModel `bun:"table:stocks"`

	Symbol    string    `bun:"symbol,type:test,pk"`
	Timestamp time.Time `bun:"timestamp,type:timestamp,pk"`
	Open      float64   `bun:"open,type:decimal,notnull"`
	Close     float64   `bun:"close,type:decimal,notnull"`
	High      float64   `bun:"high,type:decimal,notnull"`
	Low       float64   `bun:"low,type:decimal,notnull"`
}

type Stocks []Stock

func (s *Stocks) Latest() Stock {
	return slices.MaxFunc(*s, func(a, b Stock) int {
		return a.Timestamp.Compare(b.Timestamp)
	})
}

func (s *Stocks) ClosingAverage() (decimal.Decimal, error) {
	close := make([]float64, 0, len(*s))
	for _, v := range *s {
		close = append(close, v.Close)
	}
	return CalcAVG(close)
}

type StockRepository struct {
	db *bun.DB
}

func NewStockRepository(db *bun.DB) *StockRepository {
	return &StockRepository{
		db: db,
	}
}

func (r *StockRepository) Save(ctx context.Context, stocks []Stock) error {
	if len(stocks) == 0 {
		return nil
	}
	_, err := r.db.NewInsert().
		Model(&stocks).
		On("CONFLICT (symbol, timestamp) DO NOTHING").
		Exec(ctx)
	return err
}
