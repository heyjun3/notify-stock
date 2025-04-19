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

func TestSaveNotifications(t *testing.T) {
	repo := notify.NewNotificationRepository(db)
	tests := []struct {
		name          string
		notifications []notify.Notification
		err           error
	}{
		{
			name:          "Save notification with empty",
			notifications: []notify.Notification{},
			err:           nil,
		},
		{
			name: "Save notification with valid data",
			notifications: []notify.Notification{
				NoError(notify.NewNotification(nil, newSymbol("N225"), "test@example.com", time.Now())),
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Save(context.Background(), tt.notifications)
			assert.NoError(t, err)
		})
	}
}

func TestGetByID(t *testing.T) {
	repo := notify.NewNotificationRepository(db)
	t.Run("notification get by id", func(t *testing.T) {
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

	t.Run("notification get by id not found", func(t *testing.T) {
		notification, err := repo.GetByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.Nil(t, notification)
	})
}
