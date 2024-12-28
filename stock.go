package notifystock

import (
	"time"

	"github.com/uptrace/bun"
)

type Stock struct {
	bun.BaseModel

	Symbol    string    `bun:"symbol,pk"`
	Timestamp time.Time `bun:"timestamp,timestampz,pk"`
	Open      float64   `bun:"open"`
	Close     float64   `bun:"close"`
	High      float64   `bun:"high"`
	Low       float64   `bun:"low"`
}
