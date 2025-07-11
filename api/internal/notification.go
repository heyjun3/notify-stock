package notifystock

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Notification struct {
	bun.BaseModel `bun:"table:notifications"`

	ID       uuid.UUID  `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	MemberID uuid.UUID  `bun:"member_id,type:uuid"`
	Time     TimeOfHour `bun:"embed:"`

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

func NewNotification(ID *uuid.UUID, memberID uuid.UUID, symbols []string, hour time.Time) (*Notification, error) {
	if ID == nil {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		ID = &id
	}
	targets := make([]*NotificationTarget, 0, len(symbols))
	for _, symbol := range symbols {
		target, err := NewNotificationTarget(nil, *ID, symbol)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}
	return &Notification{
		ID:       *ID,
		MemberID: memberID,
		Time:     NewTimeOfHour(hour),
		Targets:  targets,
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
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := r.db.NewInsert().
			Model(&n).
			On("CONFLICT (id) DO UPDATE").
			Set(strings.Join([]string{
				"member_id = EXCLUDED.member_id",
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
	})
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

func (r *NotificationRepository) GetByMemberID(ctx context.Context, memberID uuid.UUID) ([]*Notification, error) {
	var n []*Notification
	err := r.db.NewSelect().
		Model(&n).
		Where("member_id = ?", memberID).
		Relation("Targets").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (r *NotificationRepository) GetByHour(ctx context.Context, time TimeOfHour) ([]Notification, error) {
	var n []Notification
	err := r.db.NewSelect().
		Model(&n).
		Relation("Targets").
		Where("hour = ?", time.Hour.Format("15:04:05")).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (r *NotificationRepository) DeleteByMemberID(ctx context.Context, memberID uuid.UUID) ([]*Notification, error) {
	var notifications []*Notification
	_, err := r.db.NewDelete().
		Model(&notifications).
		Where("member_id = ?", memberID).
		Returning("*").
		Exec(ctx)
	return notifications, err
}

type NotificationCreator struct {
	notificationRepository *NotificationRepository
	symbolRepository       *SymbolRepository
}

func NewNotificationCreator(
	notificationRepository *NotificationRepository,
	symbolRepository *SymbolRepository,
) *NotificationCreator {
	return &NotificationCreator{
		notificationRepository: notificationRepository,
		symbolRepository:       symbolRepository,
	}
}

func (n *NotificationCreator) Create(
	ctx context.Context, memberID uuid.UUID, symbols []string, hour time.Time,
) (*Notification, error) {
	symbolDetails, err := n.symbolRepository.GetBySymbols(ctx, symbols)
	if err != nil {
		return nil, err
	}
	if len(symbolDetails) != len(symbols) {
		return nil, fmt.Errorf("unsupported symbols: %v", symbols)
	}
	notification, err := NewNotification(nil, memberID, symbols, hour)
	if err != nil {
		return nil, err
	}
	if _, err := n.notificationRepository.DeleteByMemberID(ctx, memberID); err != nil {
		return nil, err
	}
	if err := n.notificationRepository.Save(ctx, []Notification{*notification}); err != nil {
		return nil, err
	}
	return notification, nil
}
