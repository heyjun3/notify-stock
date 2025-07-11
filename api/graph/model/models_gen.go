// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Node interface {
	IsNode()
	GetID() string
}

type ChartInput struct {
	Symbol *string   `json:"symbol,omitempty"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}

type Mutation struct {
}

type Notification struct {
	ID      string          `json:"id"`
	Time    time.Time       `json:"time"`
	Targets []*SymbolDetail `json:"targets"`
}

func (Notification) IsNode()            {}
func (this Notification) GetID() string { return this.ID }

type NotificationInput struct {
	Symbols []string  `json:"symbols"`
	Time    time.Time `json:"time"`
}

type Query struct {
}

type Stock struct {
	Symbol    string  `json:"symbol"`
	Timestamp string  `json:"timestamp"`
	Price     float64 `json:"price"`
}

type Symbol struct {
	ID     string        `json:"id"`
	Symbol string        `json:"symbol"`
	Detail *SymbolDetail `json:"detail"`
	Chart  []*Stock      `json:"chart"`
}

func (Symbol) IsNode()            {}
func (this Symbol) GetID() string { return this.ID }

type SymbolDetail struct {
	ID             string  `json:"id"`
	Symbol         string  `json:"symbol"`
	ShortName      string  `json:"shortName"`
	LongName       string  `json:"longName"`
	Price          float64 `json:"price"`
	Change         string  `json:"change"`
	ChangePercent  string  `json:"changePercent"`
	Volume         *string `json:"volume,omitempty"`
	MarketCap      *string `json:"marketCap,omitempty"`
	CurrencySymbol string  `json:"currencySymbol"`
}

type SymbolInput struct {
	Symbol string `json:"symbol"`
}
