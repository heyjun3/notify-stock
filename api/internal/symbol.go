package notifystock

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type BiMap[T comparable, U comparable] struct {
	forward  map[T]U
	backward map[U]T
}

func NewBiMap[T comparable, U comparable]() *BiMap[T, U] {
	return &BiMap[T, U]{
		forward:  make(map[T]U),
		backward: make(map[U]T),
	}
}
func NewBiMapFromMap[T comparable, U comparable](m map[T]U) *BiMap[T, U] {
	biMap := NewBiMap[T, U]()
	for k, v := range m {
		biMap.Insert(k, v)
	}
	return biMap
}

func (b *BiMap[T, U]) Insert(key T, value U) {
	if _, ok := b.forward[key]; ok {
		delete(b.backward, b.forward[key])
	}
	b.forward[key] = value
	b.backward[value] = key
}
func (b *BiMap[T, U]) Get(key T) (U, bool) {
	v, ok := b.forward[key]
	return v, ok
}
func (b *BiMap[T, U]) GetBackward(key U) (T, bool) {
	v, ok := b.backward[key]
	return v, ok
}

const (
	N225  = "N225"
	SP500 = "S&P500"
)

var (
	symbolMapForFinance = NewBiMapFromMap(map[string]string{
		N225:  "^N225",
		SP500: "^GSPC",
	})
	symbolMap = NewBiMapFromMap(map[string]string{
		"N225": "N225",
		SP500:  "S&P500",
	})
	display = NewBiMapFromMap(map[string]string{
		"N225": "Nikkei 225",
		SP500:  "S&P 500",
	})
	symbolMaps = []*BiMap[string, string]{
		symbolMapForFinance,
		symbolMap,
		display,
	}
)

type Symbol struct {
	symbol string
}

func NewSymbol(symbol string) (Symbol, error) {
	for _, m := range symbolMaps {
		value, ok := m.GetBackward(symbol)
		if ok {
			return Symbol{
				symbol: value,
			}, nil
		}
	}
	return Symbol{}, fmt.Errorf("unsupported symbol value: %s", symbol)
}

func (s Symbol) ForFinance() (string, error) {
	v, ok := symbolMapForFinance.Get(s.symbol)
	if !ok {
		return "", fmt.Errorf("unsupported finance symbol value: %s", s.symbol)
	}
	return v, nil
}

func (s Symbol) ForDB() (string, error) {
	v, ok := symbolMap.Get(s.symbol)
	if !ok {
		return "", fmt.Errorf("unsupported db symbol value: %s", s.symbol)
	}
	return v, nil
}

func (s Symbol) Display() (string, bool) {
	return display.Get(s.symbol)
}

type SymbolDetail struct {
	bun.BaseModel `bun:"table:symbols"`

	Symbol        string          `bun:"symbol,pk"`
	ShortName     string          `bun:"short_name"`
	LongName      string          `bun:"long_name"`
	MarketPrice   decimal.Decimal `bun:"market_price"`
	PreviousClose decimal.Decimal `bun:"previous_close"`
	Volume        sql.NullInt64   `bun:"volume"`
	MarketCap     sql.NullInt64   `bun:"market_cap"`
}

func (s *SymbolDetail) Change() string {
	change := s.MarketPrice.Sub(s.PreviousClose)
	if change.IsPositive() {
		return "+" + change.String()
	}
	return change.String()
}
func (s *SymbolDetail) ChangePercent() string {
	p := s.MarketPrice.Sub(s.PreviousClose).Div(s.PreviousClose).
		Mul(decimal.New(100, 0)).Round(2)
	if p.IsPositive() {
		return "+" + p.String() + "%"
	}
	return p.String() + "%"
}

type SymbolDetailOption func(detail *SymbolDetail) *SymbolDetail

func WithVolume(volume int64) SymbolDetailOption {
	return func(detail *SymbolDetail) *SymbolDetail {
		detail.Volume = sql.NullInt64{Int64: volume, Valid: true}
		return detail
	}
}
func WithMarketCap(marketCap int64) SymbolDetailOption {
	return func(detail *SymbolDetail) *SymbolDetail {
		detail.MarketCap = sql.NullInt64{Int64: marketCap, Valid: true}
		return detail
	}
}

func NewSymbolDetail(symbol string, shortName string, longName string,
	marketPrice decimal.Decimal, previousClose decimal.Decimal,
	options ...SymbolDetailOption) *SymbolDetail {
	detail := &SymbolDetail{
		Symbol:        symbol,
		ShortName:     shortName,
		LongName:      longName,
		MarketPrice:   marketPrice,
		PreviousClose: previousClose,
	}
	for _, option := range options {
		option(detail)
	}
	return detail
}

type SymbolRepository struct {
	db *bun.DB
}

func NewSymbolRepository(db *bun.DB) *SymbolRepository {
	return &SymbolRepository{
		db: db,
	}
}

func (r *SymbolRepository) Save(ctx context.Context, details []SymbolDetail) error {
	_, err := r.db.NewInsert().Model(&details).
		On("CONFLICT (symbol) DO UPDATE").
		Set(strings.Join([]string{
			"short_name = EXCLUDED.short_name",
			"long_name = EXCLUDED.long_name",
			"market_price = EXCLUDED.market_price",
			"previous_close = EXCLUDED.previous_close",
			"volume = EXCLUDED.volume",
			"market_cap = EXCLUDED.market_cap",
		}, ",")).
		Exec(ctx)
	return err
}

func (r *SymbolRepository) Get(ctx context.Context, symbol string) (*SymbolDetail, error) {
	var detail SymbolDetail
	err := r.db.NewSelect().Model(&detail).
		Where("symbol = ?", symbol).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

func (r *SymbolRepository) GetAll(ctx context.Context) ([]SymbolDetail, error) {
	var details []SymbolDetail
	err := r.db.NewSelect().Model(&details).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return details, nil
}

type SymbolFetcher struct {
	symbolRepository *SymbolRepository
}

func NewSymbolFetcher(symbolRepository *SymbolRepository) *SymbolFetcher {
	return &SymbolFetcher{
		symbolRepository: symbolRepository,
	}
}
func (f *SymbolFetcher) Fetch(ctx context.Context, symbol string) (*SymbolDetail, error) {
	sym, err := f.symbolRepository.Get(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return sym, nil
}
func (f *SymbolFetcher) FetchAll(ctx context.Context) ([]SymbolDetail, error) {
	details, err := f.symbolRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return details, nil
}
