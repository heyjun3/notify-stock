package notifystock

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

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

	Symbol    string    `bun:"symbol,type:text,pk"`
	Timestamp time.Time `bun:"timestamp,type:timestamp,pk"`
	Open      float64   `bun:"open,type:decimal,notnull"`
	Close     float64   `bun:"close,type:decimal,notnull"`
	High      float64   `bun:"high,type:decimal,notnull"`
	Low       float64   `bun:"low,type:decimal,notnull"`
}

type Stocks struct {
	symbol   Symbol
	currency Currency
	stocks   []Stock
}

func NewStocks(symbol Symbol, currency string, stocks []Stock) (*Stocks, error) {
	cur, err := CurrencyString(currency)
	if err != nil {
		return nil, err
	}
	return &Stocks{
		symbol:   symbol,
		currency: cur,
		stocks:   stocks,
	}, nil
}

func (s *Stocks) Latest() Stock {
	return slices.MaxFunc(s.stocks, func(a, b Stock) int {
		return a.Timestamp.Compare(b.Timestamp)
	})
}

func (s *Stocks) ClosingAverage() (decimal.Decimal, error) {
	close := make([]float64, 0, len(s.stocks))
	for _, v := range s.stocks {
		close = append(close, v.Close)
	}
	return CalcAVG(close)
}

func (s *Stocks) ClosingPriceToAVGRatio() (decimal.Decimal, error) {
	avg, err := s.ClosingAverage()
	if err != nil {
		return decimal.Decimal{}, err
	}
	latest := s.Latest()
	ratio := decimal.NewFromFloat(latest.Close).Div(avg)
	return ratio, nil
}

func (s *Stocks) GenerateNotificationMessage() (string, error) {
	avg, err := s.ClosingAverage()
	if err != nil {
		return "", err
	}
	latest := s.Latest()
	ratio, err := s.ClosingPriceToAVGRatio()
	if err != nil {
		return "", err
	}
	currency := s.currency
	symbolStr, _ := s.symbol.Display()
	text := strings.Join([]string{
		symbolStr,
		fmt.Sprintf("Closing Price: %v %s", int(latest.Close), currency),
		fmt.Sprintf("1-Year Moving Average: %v %s", avg.Ceil(), currency),
		fmt.Sprintf("Closing Price to 1-Year Moving Average Ratio: %v%s", ratio.Mul(decimal.New(100, 0)).RoundCeil(2), "%"),
	}, "\n")
	return text, nil
}

func CalcAVG[T cmp.Ordered](values []T) (decimal.Decimal, error) {
	d := make([]decimal.Decimal, 0, len(values))
	for _, v := range values {
		deci, err := decimal.NewFromString(fmt.Sprintf("%v", v))
		if err != nil {
			logger.Error("failed convert string to decimal")
			return decimal.Decimal{}, err
		}
		d = append(d, deci)
	}
	avg := decimal.Avg(d[0], d[1:]...)
	return avg, nil
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
	[]Stock, error) {
	s, err := symbol.ForDB()
	if err != nil {
		return nil, err
	}
	var stocks []Stock
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
		Limit(0).
		Scan(ctx); err != nil {
		return nil, err
	}
	return &stock, nil
}
