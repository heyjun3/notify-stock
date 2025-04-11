package notifystock

import "github.com/uptrace/bun"

type Notification struct {
	bun.BaseModel `bun:"table:notifications"`

	Symbol string `bun:"symbol,type:text,notnull"`
	Email  string `bun:"email,type:text,notnull"`
}

func NewNotification(symbol Symbol) (*Notification, error) {
	s, err := symbol.ForDB()
	if err != nil {
		return nil, err
	}
	return &Notification{
		Symbol: s,
	}, nil
}
