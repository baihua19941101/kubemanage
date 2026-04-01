package service

import (
	"strings"
	"sync"
	"time"
)

type AuditRecord struct {
	Time       string `json:"time"`
	RequestID  string `json:"requestId,omitempty"`
	User       string `json:"user"`
	Role       string `json:"role"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Namespace  string `json:"namespace,omitempty"`
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error,omitempty"`
}

type AuditService struct {
	mu      sync.RWMutex
	records []AuditRecord
}

type AuditFilter struct {
	User       string
	Role       string
	Method     string
	Path       string
	StatusCode int
	Limit      int
}

func NewAuditService() *AuditService {
	return &AuditService{
		records: make([]AuditRecord, 0, 64),
	}
}

func (s *AuditService) Append(requestID, user, role, method, path, namespace string, statusCode int, errMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records = append(s.records, AuditRecord{
		Time:       time.Now().Format(time.RFC3339),
		RequestID:  requestID,
		User:       user,
		Role:       role,
		Method:     method,
		Path:       path,
		Namespace:  namespace,
		StatusCode: statusCode,
		Error:      errMsg,
	})
}

func (s *AuditService) List(filter AuditFilter) []AuditRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]AuditRecord, 0, len(s.records))
	for _, item := range s.records {
		if filter.User != "" && item.User != filter.User {
			continue
		}
		if filter.Role != "" && item.Role != filter.Role {
			continue
		}
		if filter.Method != "" && item.Method != filter.Method {
			continue
		}
		if filter.Path != "" && !strings.Contains(item.Path, filter.Path) {
			continue
		}
		if filter.StatusCode > 0 && item.StatusCode != filter.StatusCode {
			continue
		}
		out = append(out, item)
	}
	if filter.Limit > 0 && len(out) > filter.Limit {
		out = out[len(out)-filter.Limit:]
	}
	return out
}
