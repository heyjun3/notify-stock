package notifystock_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
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

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
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

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		notification, err := repo.GetByID(context.Background(), id)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if notification.ID != id {
			t.Errorf("expected id %v, got %v", id, notification.ID)
		}
	})

	t.Run("notification get by id not found", func(t *testing.T) {
		notification, err := repo.GetByID(context.Background(), uuid.New())
		if err == nil {
			t.Errorf("expected no error, got %v", err)
		}
		if notification != nil {
			t.Errorf("expected nil, got %v", notification)
		}
	})
}
