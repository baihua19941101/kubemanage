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

func (s *TerminalSessionStore) TTL() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ttl
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

func (s *TerminalSessionStore) Get(id, podName string) (TerminalSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	session, err := s.getLocked(id, podName, now)
	s.cleanupExpiredLocked(now)
	return session, err
}

func (s *TerminalSessionStore) Consume(id, podName string) (TerminalSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	session, err := s.getLocked(id, podName, now)
	if err != nil {
		s.cleanupExpiredLocked(now)
		return TerminalSession{}, err
	}
	delete(s.sessions, id)
	s.cleanupExpiredLocked(now)
	return session, nil
}

func (s *TerminalSessionStore) getLocked(id, podName string, now time.Time) (TerminalSession, error) {
	session, ok := s.sessions[id]
	if !ok {
		return TerminalSession{}, ErrTerminalSessionNotFound
	}
	if now.After(session.ExpiresAt) {
		return TerminalSession{}, ErrTerminalSessionExpired
	}
	if podName != "" && session.PodName != podName {
		return TerminalSession{}, ErrTerminalSessionNotFound
	}
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
