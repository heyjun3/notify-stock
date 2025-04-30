package notifystock

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Notification struct {
	bun.BaseModel `bun:"table:notifications"`

	ID     uuid.UUID  `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	Symbol string     `bun:"symbol,type:text,notnull"`
	Email  string     `bun:"email,type:text,notnull"`
	Time   TimeOfHour `bun:"embed:"`
}

type TimeOfHour struct {
	Hour time.Time `bun:"hour,type:time,notnull"`
}

func NewTimeOfHour(hour time.Time) TimeOfHour {
	rounded := hour.Round(time.Hour).UTC()
	return TimeOfHour{
		Hour: rounded,
	}
}

func NewNotification(ID *uuid.UUID, symbol string, email string, hour time.Time) (*Notification, error) {
	if ID == nil {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		ID = &id
	}
	return &Notification{
		ID:     *ID,
		Symbol: symbol,
		Email:  email,
		Time:   NewTimeOfHour(hour),
	}, nil
}

type NotificationRepository struct {
	db *bun.DB
}

func NewNotificationRepository(db *bun.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) Save(ctx context.Context, n []Notification) error {
	if len(n) == 0 {
		return nil
	}
	_, err := r.db.NewInsert().
		Model(&n).
		On("CONFLICT (id) DO UPDATE").
		Set(strings.Join([]string{
			"symbol = EXCLUDED.symbol",
			"email = EXCLUDED.email",
			"hour = EXCLUDED.hour",
		}, ",")).
		Exec(ctx)
	return err
}
func (r *NotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*Notification, error) {
	var n Notification
	err := r.db.NewSelect().
		Model(&n).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &n, nil
}
func (r *NotificationRepository) GetByHour(ctx context.Context, time TimeOfHour) ([]Notification, error) {
	var n []Notification
	err := r.db.NewSelect().
		Model(&n).
		Where("hour = ?", time.Hour.Format("15:04:05")).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return n, nil
}

type NotificationCreator struct {
	notificationRepository *NotificationRepository
}

func NewNotificationCreator(notificationRepository *NotificationRepository) *NotificationCreator {
	return &NotificationCreator{
		notificationRepository: notificationRepository,
	}
}

func (n *NotificationCreator) Create(
	ctx context.Context, symbol string, email string, hour time.Time) (
	*Notification, error) {
	notification, err := NewNotification(nil, symbol, email, hour)
	if err != nil {
		return nil, err
	}
	if err := n.notificationRepository.Save(ctx, []Notification{*notification}); err != nil {
		return nil, err
	}
	return notification, nil
}
