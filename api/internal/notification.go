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

	Targets []*NotificationTarget `bun:"rel:has-many,join:id=notification_id"`
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
	target, err := NewNotificationTarget(nil, *ID, symbol)
	if err != nil {
		return nil, err
	}
	return &Notification{
		ID:      *ID,
		Symbol:  symbol,
		Email:   email,
		Time:    NewTimeOfHour(hour),
		Targets: []*NotificationTarget{target},
	}, nil
}

type NotificationTarget struct {
	bun.BaseModel `bun:"table:notification_targets"`

	ID             uuid.UUID `bun:"id,type:uuid,pk"`
	NotificationID uuid.UUID `bun:"notification_id,type:uuid,notnull"`
	Symbol         string    `bun:"symbol,type:text,notnull"`
}

func NewNotificationTarget(ID *uuid.UUID, notificationID uuid.UUID, symbol string) (*NotificationTarget, error) {
	if ID == nil {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		ID = &id
	}
	return &NotificationTarget{
		ID:             *ID,
		NotificationID: notificationID,
		Symbol:         symbol,
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
	if err != nil {
		return err
	}
	targets := make([]*NotificationTarget, 0, len(n))
	for _, notification := range n {
		targets = append(targets, notification.Targets...)
	}
	if len(targets) == 0 {
		return nil
	}
	_, err = r.db.NewInsert().
		Model(&targets).
		On("CONFLICT (id) DO UPDATE").
		Set(strings.Join([]string{
			"notification_id = EXCLUDED.notification_id",
			"symbol = EXCLUDED.symbol",
		}, ",")).
		Exec(ctx)
	return err
}
func (r *NotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*Notification, error) {
	var n Notification
	err := r.db.NewSelect().
		Model(&n).
		Where("id = ?", id).
		Relation("Targets").
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
func (r *NotificationRepository) GetByEmail(
	ctx context.Context, email string) ([]Notification, error) {
	var n []Notification
	if err := r.db.NewSelect().
		Model(&n).
		Where("email = ?", email).
		Scan(ctx); err != nil {
		return nil, err
	}
	return n, nil
}

type NotificationFetcher struct {
	notificationRepository *NotificationRepository
}

func NewNotificationFetcher(notificationRepository *NotificationRepository) *NotificationFetcher {
	return &NotificationFetcher{
		notificationRepository: notificationRepository,
	}
}
func (n *NotificationFetcher) GetByEmail(
	ctx context.Context, email string) ([]Notification, error) {
	return n.notificationRepository.GetByEmail(ctx, email)
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
