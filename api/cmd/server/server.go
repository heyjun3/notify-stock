package server

import (
	"log"
	"log/slog"
	"net/http"
	"os"
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
	allowedOrigin := []string{
		"http://localhost:5173",
		"https://web-server-166226611413.us-west1.run.app",
		"https://marketwatcher.shop",
	}
	origins := make(map[string]struct{})
	for _, origin := range allowedOrigin {
		origins[origin] = struct{}{}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if _, ok := origins[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
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
	logger := notifystock.CreateLogger(notifystock.Cfg.LogLevel)
	logger.Info("Start set up server")

	db := notifystock.NewDB(notifystock.Cfg.DBDSN)
	sessionRepo := notifystock.NewSessionRepository(db)
	sessions := notifystock.InitSessionsWithRepo(sessionRepo)

	resolver := graph.InitResolver(db)
	directives := graph.InitRootDirective(logger)
	c := graph.Config{
		Resolvers:  resolver,
		Directives: *directives,
	}
	srv := handler.New(graph.NewExecutableSchema(c))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	mux := http.NewServeMux()
	if notifystock.Cfg.IsDevelopment() {
		mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}
	mux.Handle("/query", notifystock.SessionMiddleware(sessions)(CORSMiddleware(loggerMiddleware(logger, srv))))
	mux.HandleFunc("GET /login", notifystock.LoginHandler(sessions))
	mux.HandleFunc("GET /auth/callback", notifystock.CallbackHandler(sessions))

	s := &http.Server{
		Addr:    "0.0.0.0" + ":" + port,
		Handler: mux,
	}
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(s.ListenAndServe())
}
