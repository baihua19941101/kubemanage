package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

var (
	ErrTerminalSessionNotFound = errors.New("terminal session not found")
	ErrTerminalSessionExpired  = errors.New("terminal session expired")
)

type TerminalSession struct {
	ID        string
	PodName   string
	Container string
	User      string
	Role      string
	ExpiresAt time.Time
}

type TerminalSessionStore struct {
	mu       sync.Mutex
	ttl      time.Duration
	sessions map[string]TerminalSession
}

func NewTerminalSessionStore(ttl time.Duration) *TerminalSessionStore {
	if ttl <= 0 {
		ttl = 2 * time.Minute
	}
	return &TerminalSessionStore{
		ttl:      ttl,
		sessions: map[string]TerminalSession{},
	}
}

func (s *TerminalSessionStore) Create(podName, container, user, role string) TerminalSession {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cleanupExpiredLocked(time.Now())
	id := "ts-" + randomHex(12)
	session := TerminalSession{
		ID:        id,
		PodName:   podName,
		Container: container,
		User:      user,
		Role:      role,
		ExpiresAt: time.Now().Add(s.ttl),
	}
	s.sessions[id] = session
	return session
}

func (s *TerminalSessionStore) Consume(id, podName string) (TerminalSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	session, ok := s.sessions[id]
	if !ok {
		s.cleanupExpiredLocked(now)
		return TerminalSession{}, ErrTerminalSessionNotFound
	}
	delete(s.sessions, id)
	if now.After(session.ExpiresAt) {
		s.cleanupExpiredLocked(now)
		return TerminalSession{}, ErrTerminalSessionExpired
	}
	if podName != "" && session.PodName != podName {
		s.cleanupExpiredLocked(now)
		return TerminalSession{}, ErrTerminalSessionNotFound
	}
	s.cleanupExpiredLocked(now)
	return session, nil
}

func (s *TerminalSessionStore) cleanupExpiredLocked(now time.Time) {
	for id, session := range s.sessions {
		if now.After(session.ExpiresAt) {
			delete(s.sessions, id)
		}
	}
}

func randomHex(size int) string {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return time.Now().Format("20060102150405")
	}
	return hex.EncodeToString(buf)
}
