package notifystock

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/uptrace/bun"
)

func InitSessionsWithRepo(repo SessionRepository) *Sessions {
	domain := "localhost"
	if Cfg.IsProduction() {
		domain = "" // 本番環境では空文字でサブドメインも含める
	}
	return NewSessions(SessionOption{
		Repository: repo,
		Domain:     domain,
		Expire:     24 * time.Hour,
	})
}

const CookieName = "notify-stock"

type SessionRepository interface {
	Get(ctx context.Context, sessionID string) (*Session, error)
	Create(ctx context.Context, session *Session) error
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, sessionID string) error
	CleanExpired(ctx context.Context) error
}

type Session struct {
	ID        string    `bun:"id,pk"`
	State     string    `bun:"state,notnull"`
	IsActive  bool      `bun:"is_active,notnull,default:false"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	ExpiresAt time.Time `bun:"expires_at,notnull"`
}

func NewSession(expire time.Duration) (*Session, error) {
	id, err := randomString()
	if err != nil {
		return nil, err
	}
	state, err := randomString()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &Session{
		ID:        id,
		State:     state,
		IsActive:  false,
		CreatedAt: now,
		ExpiresAt: now.Add(expire),
	}, nil
}

type SessionDatabaseRepository struct {
	db *bun.DB
}

func NewSessionRepository(db *bun.DB) SessionRepository {
	return &SessionDatabaseRepository{db: db}
}

func (r *SessionDatabaseRepository) Get(ctx context.Context, sessionID string) (*Session, error) {
	session := &Session{}
	err := r.db.NewSelect().
		Model(session).
		Where("id = ? AND expires_at > NOW()", sessionID).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, err
	}
	return session, nil
}

func (r *SessionDatabaseRepository) Create(ctx context.Context, session *Session) error {
	_, err := r.db.NewInsert().
		Model(session).
		Exec(ctx)
	return err
}

func (r *SessionDatabaseRepository) Update(ctx context.Context, session *Session) error {
	_, err := r.db.NewUpdate().
		Model(session).
		Where("id = ?", session.ID).
		Exec(ctx)
	return err
}

func (r *SessionDatabaseRepository) Delete(ctx context.Context, sessionID string) error {
	_, err := r.db.NewDelete().
		Model((*Session)(nil)).
		Where("id = ?", sessionID).
		Exec(ctx)
	return err
}

func (r *SessionDatabaseRepository) CleanExpired(ctx context.Context) error {
	_, err := r.db.NewDelete().
		Model((*Session)(nil)).
		Where("expires_at <= NOW()").
		Exec(ctx)
	return err
}

type Sessions struct {
	repo   SessionRepository
	domain string
	expire time.Duration
}
type SessionOption struct {
	Repository SessionRepository
	Domain     string
	Expire     time.Duration
}

func NewSessions(opt SessionOption) *Sessions {
	return &Sessions{
		repo:   opt.Repository,
		domain: opt.Domain,
		expire: opt.Expire,
	}
}

func NewSessionsWithDefaults(repo SessionRepository) *Sessions {
	return &Sessions{
		repo:   repo,
		domain: "localhost",
		expire: 24 * time.Hour,
	}
}

func (s *Sessions) Get(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return nil, err
	}
	if cookie == nil {
		return nil, fmt.Errorf("not found cookie")
	}

	return s.repo.Get(r.Context(), cookie.Value)
}

func (s *Sessions) New(w http.ResponseWriter, ctx context.Context) (*Session, error) {
	session, err := NewSession(s.expire)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(ctx, session)
	if err != nil {
		return nil, err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return session, nil
}

func (s *Sessions) Store(ctx context.Context, session *Session) error {
	return s.repo.Update(ctx, session)
}

func (s *Sessions) CleanExpired(ctx context.Context) error {
	return s.repo.CleanExpired(ctx)
}

func (s *Sessions) Delete(ctx context.Context, sessionID string) error {
	return s.repo.Delete(ctx, sessionID)
}

func randomString() (string, error) {
	randBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, randBytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(randBytes), nil
}

type sessionKeyType string

var sessionKey = sessionKeyType("session")

func SessionMiddleware(sessions *Sessions) func(next http.Handler) http.Handler {
	return (func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := sessions.Get(r)
			if err != nil {
				logger.Info(err.Error())
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), sessionKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

func GetSession(ctx context.Context) (*Session, error) {
	session, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return session, nil
}
