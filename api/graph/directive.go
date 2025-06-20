package graph

import (
	"context"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	notifystock "github.com/heyjun3/notify-stock/internal"
)

type Directive func(
	ctx context.Context, obj any, next graphql.Resolver,
) (any, error)

func NewAuthDirective(logger *slog.Logger) Directive {
	return func(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
		_, err := notifystock.GetSession(ctx)
		if err != nil {
			logger.Error("failed to get session")
			return nil, &gqlerror.Error{
				Err:     err,
				Message: err.Error(),
				Extensions: map[string]any{
					"code": "UNAUTHORIZED",
				},
			}
		}
		return next(ctx)
	}
}

func NewDirectiveRoot(auth Directive) *DirectiveRoot {
	return &DirectiveRoot{
		Auth: auth,
	}
}
