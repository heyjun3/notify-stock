package notifystock

import "fmt"

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
