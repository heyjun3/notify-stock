package notifystock_test

import (
	"context"
	"testing"
	"time"

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
