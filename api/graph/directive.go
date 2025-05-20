package graph

import (
	"context"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
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
			return nil, err
		}
		return next(ctx)
	}
}

func NewDirectiveRoot(auth Directive) *DirectiveRoot {
	return &DirectiveRoot{
		Auth: auth,
	}
}
