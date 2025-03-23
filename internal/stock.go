package notifystock

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

func NewStock(symbol Symbol, timestamp time.Time, currency string,
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
	cur, err := CurrencyString(currency)
	if err != nil {
		return Stock{}, err
	}
	return Stock{
		symbol:    symbol,
		currency:  cur,
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

	symbol   Symbol   `bun:"-"`
	currency Currency `bun:"-"`

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

func (s *Stocks) GenerateNotificationMessage() (string, error) {
	avg, err := s.ClosingAverage()
	if err != nil {
		return "", err
	}
	latest := s.Latest()
	currency := latest.currency.String()
	symbolStr, _ := latest.symbol.Display()
	text := strings.Join([]string{
		symbolStr,
		fmt.Sprintf("Closing Price: %v %s", int(latest.Close), currency),
		fmt.Sprintf("1-Year Moving Average: %v %s", avg.Ceil(), currency),
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
