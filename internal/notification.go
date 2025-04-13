package notifystock

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Notification struct {
	bun.BaseModel `bun:"table:notifications"`

	ID     uuid.UUID `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	Symbol string    `bun:"symbol,type:text,notnull"`
	Email  string    `bun:"email,type:text,notnull"`
	Hour TimeOfHour `bun:"embed:"`
}

type TimeOfHour struct {
	Hour time.Time `bun:"hour,type:time,notnull"`
}
func NewTimeOfHour(hour time.Time) TimeOfHour {
	rounded := hour.Round(time.Hour)
	return TimeOfHour{
		Hour: rounded,
	}
}

func NewNotification(ID *uuid.UUID, symbol Symbol, email string, hour time.Time) (*Notification, error) {
	if ID == nil {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		ID = &id
	}
	s, err := symbol.ForDB()
	if err != nil {
		return nil, err
	}
	return &Notification{
		ID:    *ID,
		Symbol: s,
		Email: email,
		Hour:  NewTimeOfHour(hour),
	}, nil
}
