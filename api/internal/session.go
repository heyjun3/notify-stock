package notifystock

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"slices"
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
	cookie := r.Header.Get("Cookie")
	if cookie == "" {
		return nil, fmt.Errorf("not found cookie header")
	}
	cookies, err := http.ParseCookie(cookie)
	if err != nil {
		return nil, err
	}
	i := slices.IndexFunc(cookies, func(c *http.Cookie) bool {
		return c.Name == CookieName
	})
	if i == -1 {
		return nil, fmt.Errorf("not found support cookie")
	}
	if session, ok := s.store[cookies[i].Value]; ok {
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
