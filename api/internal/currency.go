package notifystock

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

//go:generate enumer -type=Currency
type Currency int

const (
	_ Currency = iota
	JPY
	USD
)

var _ driver.Valuer = (*Currency)(nil)

func (c Currency) Value() (driver.Value, error) {
	if c.IsACurrency() {
		return c.String(), nil
	}
	return nil, nil
}

var _ sql.Scanner = (*Currency)(nil)

func (c *Currency) Scan(value any) (err error) {
	switch v := value.(type) {
	case string:
		*c, err = CurrencyString(v)
		return err
	case nil:
		c = nil
		return
	default:
		return fmt.Errorf("unsupported type %T for Currency", value)
	}
}
