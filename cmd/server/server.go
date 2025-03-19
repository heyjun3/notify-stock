package server

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/spf13/cobra"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/heyjun3/notify-stock/graph"
	notifystock "github.com/heyjun3/notify-stock/internal"
)

var ServerCommand = &cobra.Command{
	Use:   "server",
	Short: "Run Server",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

const defaultPort = "8080"

func loggerMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("request start")
		next.ServeHTTP(w, r)
		logger.Info("request end")
	})
}

func runServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	repo := notifystock.InitStockRepository(notifystock.Cfg.DBDSN)
	resolver := graph.NewResolver(repo)
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", loggerMiddleware(logger, srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
