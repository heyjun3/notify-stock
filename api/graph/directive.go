package graph

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/vektah/gqlparser/v2/gqlerror"

	notifystock "github.com/heyjun3/notify-stock/internal"
)

type memberIDKeyType struct{}

var memberIDKey = memberIDKeyType{}

func GetMemberID(ctx context.Context) (*uuid.UUID, error) {
	if id, ok := ctx.Value(memberIDKey).(*uuid.UUID); ok {
		return id, nil
	}
	return nil, fmt.Errorf("not found memberID")
}
func SetMemberID(ctx context.Context, memberID uuid.UUID) context.Context {
	return context.WithValue(ctx, memberIDKey, &memberID)
}

type Directive func(
	ctx context.Context, obj any, next graphql.Resolver,
) (any, error)

func NewAuthDirective(logger *slog.Logger) Directive {
	return func(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
		session, err := notifystock.GetSession(ctx)
		if err != nil {
			logger.Warn("failed to get session")
			return nil, &gqlerror.Error{
				Err:     err,
				Message: err.Error(),
				Extensions: map[string]any{
					"code": "UNAUTHORIZED",
				},
			}
		}
		c := SetMemberID(ctx, session.MemberID)
		return next(c)
	}
}

func NewDirectiveRoot(auth Directive) *DirectiveRoot {
	return &DirectiveRoot{
		Auth: auth,
	}
}
