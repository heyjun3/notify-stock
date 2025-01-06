package notifystock_test

import (
	"context"
	"testing"
	"time"

	notify "github.com/heyjun3/notify-stock"
	"github.com/stretchr/testify/assert"

	"github.com/uptrace/bun"
)

var db *bun.DB

func init() {
	dsn := "postgres://postgres:postgres@localhost:5555/notify-stock-test?sslmode=disable"
	db = notify.NewDB(dsn)
}

func TestSave(t *testing.T) {
	repo := notify.NewStockRepository(db)

	t.Run("", func(t *testing.T) {
		stocks := []notify.Stock{
			{Symbol: "N255", Timestamp: time.Now(), Open: 1000, Close: 2000, High: 2500, Low: 500},
		}
		err := repo.Save(context.Background(), stocks)

		assert.NoError(t, err)
	})
}
