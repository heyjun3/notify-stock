package notifystock_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	notify "github.com/heyjun3/notify-stock/internal"
	"github.com/stretchr/testify/assert"
)

func TestMemberRepository(t *testing.T) {
	db := openDB(t)
	repo := notify.NewMemberRepository(db)

	t.Run("save and get member", func(t *testing.T) {
		ctx := context.Background()
		v7, err := uuid.NewV7()
		assert.NoError(t, err)

		member, err := notify.NewGoogleMember(
			&v7,
			"google-id",
			"email",
			true,
			"Name",
			"GivenName",
			"FamilyName",
			"PictureURL",
		)
		assert.NoError(t, err)

		err = repo.Save(ctx, []*notify.Member{member})
		assert.NoError(t, err)

		savedMember, err := repo.GetByID(ctx, member.ID)
		assert.NoError(t, err)
		assert.Equal(t, member, savedMember)
		assert.Equal(t, member.GoogleMember, savedMember.GoogleMember)

		googleMember, err := repo.GetByGoogleID(ctx, "google-id")
		assert.NoError(t, err)
		assert.Equal(t, member, googleMember)
		assert.Equal(t, member.GoogleMember, googleMember.GoogleMember)
	})

	t.Run("get or create google member", func(t *testing.T) {
		ctx := context.Background()
		member, err := notify.NewGoogleMember(
			nil,
			"google-id-1",
			"email",
			true,
			"Name",
			"GivenName",
			"FamilyName",
			"PictureURL",
		)
		assert.NoError(t, err)

		createdMember, err := repo.GetOrCreateGoogleMember(ctx, member)
		assert.NoError(t, err)
		assert.Equal(t, member, createdMember)

		existsMember, err := repo.GetOrCreateGoogleMember(ctx, member)
		assert.NoError(t, err)
		assert.Equal(t, member, existsMember)

		invalidMember, err := notify.NewMember(nil)
		assert.NoError(t, err)

		_, err = repo.GetOrCreateGoogleMember(ctx, invalidMember)
		assert.Error(t, err, "member does not have a GoogleMember")
	})

	t.Run("not found errors", func(t *testing.T) {
		ctx := context.Background()
		v7, err := uuid.NewV7()
		assert.NoError(t, err)

		_, err = repo.GetByID(ctx, v7)
		assert.Error(t, err)

		_, err = repo.GetByGoogleID(ctx, "non-existent-google-id")
		assert.Error(t, err)
	})

	t.Run("save invalid member", func(t *testing.T) {
		ctx := context.Background()
		invalidMember, err := notify.NewMember(nil)
		assert.NoError(t, err)
		invalidMember.GoogleMember = &notify.GoogleMember{}

		err = repo.Save(ctx, []*notify.Member{invalidMember})
		assert.Error(t, err)
	})

	t.Run("get or create google member with nil google member", func(t *testing.T) {
		ctx := context.Background()
		member, err := notify.NewMember(nil)
		assert.NoError(t, err)

		_, err = repo.GetOrCreateGoogleMember(ctx, member)
		assert.Error(t, err)
	})

	t.Run("save empty slice", func(t *testing.T) {
		ctx := context.Background()
		err := repo.Save(ctx, []*notify.Member{})
		assert.NoError(t, err)
	})

	t.Run("save member with no google member", func(t *testing.T) {
		ctx := context.Background()
		member, err := notify.NewMember(nil)
		assert.NoError(t, err)

		err = repo.Save(ctx, []*notify.Member{member})
		assert.NoError(t, err)
	})
}
