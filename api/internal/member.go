package notifystock

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Member struct {
	bun.BaseModel `bun:"table:members"`

	ID uuid.UUID `bun:"id,type:uuid,pk"`

	GoogleMember *GoogleMember `bun:"rel:has-one,join:id=member_id"`
}

func NewMember(id *uuid.UUID) (*Member, error) {
	if id == nil {
		i, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		id = &i
	}
	return &Member{
		ID: *id,
	}, nil
}

type GoogleMember struct {
	bun.BaseModel `bun:"table:google_members"`

	ID            string    `bun:"id,type:text,pk"`
	Email         string    `bun:"email,type:text,notnull"`
	VerifiedEmail bool      `bun:"verified_email,type:boolean,notnull"`
	Name          string    `bun:"name,type:text,notnull"`
	GivenName     string    `bun:"given_name,type:text,notnull"`
	FamilyName    string    `bun:"family_name,type:text,notnull"`
	Picture       string    `bun:"picture,type:text,notnull"`
	MemberID      uuid.UUID `bun:"member_id,type:uuid,notnull"`
}

func NewGoogleMember(
	ID *uuid.UUID,
	googleID string,
	email string,
	verifiedEmail bool,
	name string,
	givenName string,
	familyName string,
	picture string,
) (*Member, error) {
	member, err := NewMember(ID)
	if err != nil {
		return nil, err
	}
	googleMember := &GoogleMember{
		ID:            googleID,
		Email:         email,
		VerifiedEmail: verifiedEmail,
		Name:          name,
		GivenName:     givenName,
		FamilyName:    familyName,
		Picture:       picture,
		MemberID:      member.ID,
	}
	member.GoogleMember = googleMember
	return member, nil
}

type MemberRepository struct {
	db *bun.DB
}

func NewMemberRepository(db *bun.DB) *MemberRepository {
	return &MemberRepository{
		db: db,
	}
}

func (r *MemberRepository) GetByID(ctx context.Context, id uuid.UUID) (*Member, error) {
	var member Member
	if err := r.db.NewSelect().
		Model(&member).
		Where("member.id = ?", id).
		Relation("GoogleMember").
		Scan(ctx); err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepository) GetByGoogleID(ctx context.Context, googleID string) (*Member, error) {
	var member Member
	if err := r.db.NewSelect().
		Model(&member).
		Relation("GoogleMember").
		Where("google_member.id = ?", googleID).
		Scan(ctx); err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepository) GetOrCreateGoogleMember(ctx context.Context, member *Member) (*Member, error) {
	if member.GoogleMember == nil {
		return nil, fmt.Errorf("member does not have a GoogleMember")
	}
	exist, err := r.db.NewSelect().
		Model((*Member)(nil)).
		Relation("GoogleMember").
		Where("google_member.id = ?", member.GoogleMember.ID).
		Exists(ctx)
	if err != nil {
		return nil, err
	}
	if exist {
		return member, nil
	}
	if err := r.Save(ctx, []*Member{member}); err != nil {
		return nil, err
	}
	return member, nil
}

func (r *MemberRepository) Save(ctx context.Context, members []*Member) error {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if len(members) == 0 {
			return nil
		}
		_, err := r.db.NewInsert().
			Model(&members).
			On("CONFLICT (id) DO NOTHING").
			Exec(ctx)
		if err != nil {
			return err
		}
		googleMembers := make([]*GoogleMember, 0, len(members))
		for _, member := range members {
			if member.GoogleMember != nil {
				googleMembers = append(googleMembers, member.GoogleMember)
			}
		}
		if len(googleMembers) == 0 {
			return nil
		}
		_, err = r.db.NewInsert().
			Model(&googleMembers).
			On("CONFLICT (id) DO UPDATE").
			Set(strings.Join([]string{
				"email = EXCLUDED.email",
				"verified_email = EXCLUDED.verified_email",
				"name = EXCLUDED.name",
				"given_name = EXCLUDED.given_name",
				"family_name = EXCLUDED.family_name",
				"picture = EXCLUDED.picture",
				"member_id = EXCLUDED.member_id",
			}, ",")).
			Exec(ctx)
		return err
	})
	return err
}
