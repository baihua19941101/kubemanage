package service

import (
	"sync"
	"time"
)

type AuditRecord struct {
	Time       string `json:"time"`
	User       string `json:"user"`
	Role       string `json:"role"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	StatusCode int    `json:"statusCode"`
}

type AuditService struct {
	mu      sync.RWMutex
	records []AuditRecord
}

func NewAuditService() *AuditService {
	return &AuditService{
		records: make([]AuditRecord, 0, 64),
	}
}

func (s *AuditService) Append(user, role, method, path string, statusCode int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records = append(s.records, AuditRecord{
		Time:       time.Now().Format(time.RFC3339),
		User:       user,
		Role:       role,
		Method:     method,
		Path:       path,
		StatusCode: statusCode,
	})
}

func (s *AuditService) List() []AuditRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]AuditRecord, len(s.records))
	copy(out, s.records)
	return out
}
