package notifystock

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"
)

var sessions *Sessions

func init() {
	sessions = NewSessions(SessionOption{
		Domain: "localhost",
		Expire: 24 * time.Hour,
	})
}

const CookieName = "notify-stock"

type Session struct {
	ID       string
	state    string
	isActive bool
}

func NewSession() (*Session, error) {
	id, err := randomString()
	if err != nil {
		return nil, err
	}
	state, err := randomString()
	if err != nil {
		return nil, err
	}
	return &Session{
		ID:       id,
		state:    state,
		isActive: false,
	}, nil
}

type Sessions struct {
	store  map[string]*Session
	domain string
	expire time.Duration
}
type SessionOption struct {
	Domain string
	Expire time.Duration
}

func NewSessions(opt SessionOption) *Sessions {
	return &Sessions{
		store:  make(map[string]*Session),
		domain: opt.Domain,
		expire: opt.Expire,
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
	if session, ok := s.store[cookie.Value]; ok {
		return session, nil
	} else {
		return nil, fmt.Errorf("session not found")
	}
}

func (s *Sessions) New(w http.ResponseWriter) (*Session, error) {
	session, err := NewSession()
	if err != nil {
		return nil, err
	}
	s.store[session.ID] = session
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    session.ID,
		Domain:   s.domain,
		Expires:  time.Now().Add(s.expire),
		HttpOnly: true,
	})
	return session, nil
}

func (s *Sessions) Store(session *Session) {
	s.store[session.ID] = session
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

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := sessions.Get(r)
		if err != nil {
			logger.Error(err.Error())
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), sessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetSession(ctx context.Context) (*Session, error) {
	session, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return session, nil
}
