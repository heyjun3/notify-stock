package notifystock_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	notify "github.com/heyjun3/notify-stock/internal"
)

func NoError[T any](t *T, err error) T {
	return *t
}

func TestNotificationRepository(t *testing.T) {
	db := openDB(t)
	repo := notify.NewNotificationRepository(db)
	t.Run("save notification with empty", func(t *testing.T) {
		err := repo.Save(context.Background(), []notify.Notification{})
		assert.NoError(t, err)
	})
	t.Run("save notification with valid data", func(t *testing.T) {
		err := repo.Save(context.Background(), []notify.Notification{
			NoError(notify.NewNotification(nil, newSymbol("N225"), "test@example.com", time.Now())),
		})
		assert.NoError(t, err)
	})
	t.Run("save notification update existing", func(t *testing.T) {
		id := uuid.New()
		ctx := context.Background()
		err := repo.Save(ctx, []notify.Notification{
			NoError(notify.NewNotification(
				&id, newSymbol("N225"), "test@example.com",
				time.Date(2022, 1, 1, 4, 0, 0, 0, time.UTC))),
		})
		assert.NoError(t, err)

		err = repo.Save(ctx, []notify.Notification{
			NoError(notify.NewNotification(
				&id, newSymbol("N225"), "test+test@example.com",
				time.Date(2022, 1, 1, 23, 0, 0, 0, time.UTC))),
		})
		assert.NoError(t, err)

		notification, err := repo.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, id, notification.ID)
		assert.Equal(t, "test+test@example.com", notification.Email)
		assert.Equal(t, 23, notification.Time.Hour.Hour())
	})
	t.Run("get notification by id", func(t *testing.T) {
		id := uuid.New()

		notifications := []notify.Notification{
			NoError(notify.NewNotification(&id, newSymbol("N225"), "test@exsample.com", time.Now())),
		}
		err := repo.Save(context.Background(), notifications)
		assert.NoError(t, err)

		notification, err := repo.GetByID(context.Background(), id)
		assert.NoError(t, err)

		assert.Equal(t, notification.ID, id)
		assert.Equal(t, notification.Symbol, notifications[0].Symbol)
		assert.Equal(t, notification.Email, notifications[0].Email)
		assert.Equal(t, notification.Time.Hour.Hour(), notifications[0].Time.Hour.Hour())
	})
	t.Run("get notification by id not found", func(t *testing.T) {
		notification, err := repo.GetByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.Nil(t, notification)
	})
	t.Run("get notifications by hour", func(t *testing.T) {
		notifications := []notify.Notification{
			NoError(notify.NewNotification(
				nil, newSymbol("N225"), "test@exsample.com", time.Now())),
			NoError(notify.NewNotification(
				nil, newSymbol("N225"), "test@exsample.com", time.Now().Add(2*time.Hour))),
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
