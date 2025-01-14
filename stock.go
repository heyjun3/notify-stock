package notifystock

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

var StockSymbolMap = map[string]string{
	N225:  "N225",
	SP500: "S&P500",
}

func NewStock(symbol string, timestamp time.Time,
	open, close, high, low float64) (Stock, error) {
	s := StockSymbolMap[symbol]
	if s == "" {
		return Stock{}, fmt.Errorf("undefined symbol error")
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
