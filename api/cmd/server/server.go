package server

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

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
		now := time.Now()
		logger.Info("request start")
		next.ServeHTTP(w, r)
		logger.Info("request end", slog.Duration("duration", time.Duration(time.Since(now).Milliseconds())))
	})
}
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if strings.Contains(origin, "localhost") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Max-Age", "86400")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func runServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	resolver := graph.InitResolver(notifystock.Cfg.DBDSN)
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	mux := http.NewServeMux()

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", CORSMiddleware(loggerMiddleware(logger, srv)))

	s := &http.Server{
		Addr:    "0.0.0.0" + ":" + port,
		Handler: mux,
	}
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(s.ListenAndServe())
}
