package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"voice-bridge-server/internal/config"
)

type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type Store struct {
	cfg      config.Config
	mu       sync.Mutex
	sessions map[string]Session
	users    *UserStore
}

func NewStore(cfg config.Config) (*Store, error) {
	users, err := NewUserStore(cfg.DBFile)
	if err != nil {
		return nil, err
	}
	return &Store{cfg: cfg, sessions: map[string]Session{}, users: users}, nil
}

func (s *Store) Register(username, password string, remember bool) (Session, bool, error) {
	if !validAccount(username, password) {
		return Session{}, false, fmt.Errorf("账号至少 3 位，密码至少 6 位")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Session{}, false, err
	}
	if !s.users.Create(username, string(hash)) {
		return Session{}, false, fmt.Errorf("账号已存在")
	}
	return s.newSession(username, remember), true, nil
}

func (s *Store) Login(username, password string, remember bool) (Session, bool) {
	user, ok := s.users.Find(username)
	if !ok {
		return Session{}, false
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return Session{}, false
	}
	return s.newSession(username, remember), true
}

func (s *Store) newSession(userID string, remember bool) Session {
	ttl := time.Duration(s.cfg.ShortTokenExpireHours) * time.Hour
	if remember {
		ttl = time.Duration(s.cfg.TokenExpireDays) * 24 * time.Hour
	}
	session := Session{
		Token:     randomToken(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl),
	}
	s.mu.Lock()
	s.sessions[session.Token] = session
	s.mu.Unlock()
	return session
}

func (s *Store) Verify(token string) (Session, bool) {
	if token == "" {
		return Session{}, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[token]
	if !ok {
		return Session{}, false
	}
	if time.Now().After(session.ExpiresAt) {
		delete(s.sessions, token)
		return Session{}, false
	}
	return session, true
}

func (s *Store) Logout(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

func randomToken() string {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}

func validAccount(username, password string) bool {
	return len(username) >= 3 && len(password) >= 6
}
