package notifystock

import (
	"context"
	"database/sql"
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

func (r *StockRepository) GetStockByPeriod(
	ctx context.Context, symbol Symbol, begging, end time.Time) (
	Stocks, error) {
	s, err := symbol.ForDB()
	if err != nil {
		return nil, err
	}
	var stocks Stocks
	if err := r.db.NewSelect().
		Model(&stocks).
		DistinctOn("timestamp::date").
		Where("symbol = ?", s).
		Where("timestamp::date BETWEEN ? AND ?", begging, end).
		OrderExpr("timestamp::date").
		Order("timestamp").
		Scan(ctx); err != nil {
		return nil, err
	}
	return stocks, nil
}

func (r *StockRepository) GetLatestStock(ctx context.Context, symbol Symbol) (*Stock, error) {
	s, err := symbol.ForDB()
	if err != nil {
		return nil, err
	}
	var stock Stock
	if err := r.db.NewSelect().
		Model(&stock).
		Where("symbol = ?", s).
		Order("timestamp DESC").
		Limit(1).
		Scan(ctx); err != nil {
		return nil, err
	}
	return &stock, nil
}
