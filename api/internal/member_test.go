package notifystock_test

import (
	"context"
	"testing"

	notify "github.com/heyjun3/notify-stock/internal"
	"github.com/stretchr/testify/assert"
)

func TestMemberRepository(t *testing.T) {
	db := openDB(t)
	repo := notify.NewMemberRepository(db)

	t.Run("save and get member", func(t *testing.T) {
		ctx := context.Background()
		member, err := notify.NewMember(nil)
		assert.NoError(t, err)

		member.GoogleMember = notify.NewGoogleMember(
			"google-id",
			"email",
			true,
			"Name",
			"GivenName",
			"FamilyName",
			"PictureURL",
			member.ID,
		)

		err = repo.Save(ctx, []*notify.Member{member})
		assert.NoError(t, err)

		savedMember, err := repo.GetByID(ctx, member.ID)
		assert.NoError(t, err)
		assert.Equal(t, member, savedMember)
		assert.Equal(t, member.GoogleMember, savedMember.GoogleMember)
	})
}
