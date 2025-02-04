package notifystock

import "fmt"

var (
	symbolMapForFinance = map[string]string{
		"N225":  "^N225",
		"SP500": "^GSPC",
	}
	symbolMapForDB = map[string]string{
		"N225":  "N225",
		"S&P500": "S&P500",
	}
	display = map[string]string{
		"N225":  "Nikkei 225",
		"SP500": "S&P 500",
	}
)

var validSymbol = map[string]string{}

func init() {
	for k, v := range symbolMapForFinance {
		validSymbol[v] = k
	}
	for k, v := range symbolMapForDB {
		validSymbol[v] = k
	}
}

type Symbol struct {
	symbol string
}

func NewSymbol(symbol string) (Symbol, error) {
	value, ok := validSymbol[symbol]
	if !ok {
		return Symbol{}, fmt.Errorf("unsupported symbol value: %s", symbol)
	}
	return Symbol{
		symbol: value,
	}, nil
}

func (s Symbol) ForFinance() (string, error) {
	v, ok := symbolMapForFinance[s.symbol]
	if !ok {
		return "", fmt.Errorf("unsupported finance symbol value: %s", s.symbol)
	}
	return v, nil
}

func (s Symbol) ForDB() (string, error) {
	v, ok := symbolMapForDB[s.symbol]
	if !ok {
		return "", fmt.Errorf("unsupported db symbol value: %s", s.symbol)
	}
	return v, nil
}

func (s Symbol) Display() string {
	return display[s.symbol]
}
