package grapherror

import (
	"context"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func NewPresenter(logger *slog.Logger) graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) *gqlerror.Error {
		logger.Error("GraphQL error", "error", err)
		e := graphql.DefaultErrorPresenter(ctx, err)
		if e.Extensions == nil {
			e.Extensions = make(map[string]any)
		}
		if _, ok := e.Extensions["code"]; !ok {
			e.Extensions["code"] = InternalServerError
		}
		return e
	}
}
