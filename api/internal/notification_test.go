package notifystock_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	notify "github.com/heyjun3/notify-stock/internal"
)

func NoError[T any](t *T, err error) T {
	return *t
}

func TestNotificationRepository(t *testing.T) {
	ctx := context.Background()
	db := openDB(t)
	member, err := notify.NewMember(nil)
	assert.NoError(t, err)
	memberRepository := notify.NewMemberRepository(db)
	if err := memberRepository.Save(ctx, []*notify.Member{member}); err != nil {
		panic(err)
	}

	symbol := notify.NewSymbolDetail("N225", "NIKKEI 225", "NIKKEI", "JPY", decimal.NewFromInt(1000), decimal.NewFromInt(1000))
	symbolRepo := notify.NewSymbolRepository(db)
	if err := symbolRepo.Save(context.Background(), []notify.SymbolDetail{*symbol}); err != nil {
		panic(err)
	}
	repo := notify.NewNotificationRepository(db)

	t.Run("save notification with empty", func(t *testing.T) {
		err := repo.Save(context.Background(), []notify.Notification{})
		assert.NoError(t, err)
	})
	t.Run("save notification with valid data", func(t *testing.T) {
		err := repo.Save(context.Background(), []notify.Notification{
			NoError(notify.NewNotification(nil, member.ID, []string{"N225"}, time.Now())),
		})
		assert.NoError(t, err)
	})
	t.Run("save notification update existing", func(t *testing.T) {
		id := uuid.New()
		ctx := context.Background()
		err := repo.Save(ctx, []notify.Notification{
			NoError(notify.NewNotification(
				&id, member.ID, []string{"N225"},
				time.Date(2022, 1, 1, 4, 0, 0, 0, time.UTC))),
		})
		assert.NoError(t, err)

		err = repo.Save(ctx, []notify.Notification{
			NoError(notify.NewNotification(
				&id, member.ID, []string{"N225"},
				time.Date(2022, 1, 1, 23, 0, 0, 0, time.UTC))),
		})
		assert.NoError(t, err)

		notification, err := repo.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, id, notification.ID)
		assert.Equal(t, 23, notification.Time.Hour.Hour())
	})
	t.Run("get notification by id", func(t *testing.T) {
		id := uuid.New()

		n, err := notify.NewNotification(&id, member.ID, []string{"N225"}, time.Now())
		assert.NoError(t, err)

		err = repo.Save(context.Background(), []notify.Notification{*n})
		assert.NoError(t, err)

		notification, err := repo.GetByID(context.Background(), id)
		assert.NoError(t, err)

		assert.Equal(t, notification.ID, id)
		assert.Equal(t, notification.Time.Hour.Hour(), n.Time.Hour.Hour())
		assert.Equal(t, notification.Targets, n.Targets)
	})
	t.Run("get notification by id not found", func(t *testing.T) {
		notification, err := repo.GetByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.Nil(t, notification)
	})
	t.Run("get notifications by hour", func(t *testing.T) {
		notifications := []notify.Notification{
			NoError(notify.NewNotification(
				nil, member.ID, []string{"N225"}, time.Now())),
			NoError(notify.NewNotification(
				nil, member.ID, []string{"N225"}, time.Now().Add(2*time.Hour))),
		}
		err := repo.Save(context.Background(), notifications)
		assert.NoError(t, err)

		ns, err := repo.GetByHour(
			context.Background(), notify.NewTimeOfHour(time.Now().UTC()))

		assert.NoError(t, err)
		assert.Greater(t, len(ns), 0)
		for _, n := range ns {
			assert.Equal(t, n.Time.Hour.Hour(), notify.NewTimeOfHour(time.Now()).Hour.Hour())
		}
	})
}
