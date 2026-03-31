package service

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

type Namespace struct {
	Name      string            `json:"name"`
	Status    string            `json:"status"`
	Age       string            `json:"age"`
	Labels    map[string]string `json:"labels,omitempty"`
	CreatedAt time.Time         `json:"-"`
}

type NamespaceService struct {
	namespaces []Namespace
}

func NewNamespaceService() *NamespaceService {
	now := time.Now()
	return &NamespaceService{
		namespaces: []Namespace{
			{Name: "default", Status: "Active", CreatedAt: now.Add(-240 * time.Hour), Labels: map[string]string{"system": "true"}},
			{Name: "kube-system", Status: "Active", CreatedAt: now.Add(-240 * time.Hour), Labels: map[string]string{"system": "true"}},
			{Name: "dev", Status: "Active", CreatedAt: now.Add(-24 * time.Hour), Labels: map[string]string{"env": "dev"}},
		},
	}
}

func (s *NamespaceService) List() []Namespace {
	items := slices.Clone(s.namespaces)
	for i := range items {
		items[i].Age = humanAge(items[i].CreatedAt)
	}
	return items
}

func (s *NamespaceService) Create(name string, labels map[string]string) (Namespace, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Namespace{}, errors.New("namespace name is required")
	}

	if _, exists := s.find(name); exists {
		return Namespace{}, fmt.Errorf("namespace already exists: %s", name)
	}

	ns := Namespace{
		Name:      name,
		Status:    "Active",
		CreatedAt: time.Now(),
		Labels:    labels,
	}
	s.namespaces = append(s.namespaces, ns)
	ns.Age = humanAge(ns.CreatedAt)
	return ns, nil
}

func (s *NamespaceService) Delete(name string) error {
	idx := -1
	for i := range s.namespaces {
		if s.namespaces[i].Name == name {
			idx = i
			break
		}
	}
	if idx < 0 {
		return fmt.Errorf("namespace not found: %s", name)
	}
	s.namespaces = append(s.namespaces[:idx], s.namespaces[idx+1:]...)
	return nil
}

func (s *NamespaceService) YAML(name string) (string, error) {
	ns, ok := s.find(name)
	if !ok {
		return "", fmt.Errorf("namespace not found: %s", name)
	}

	labels := ""
	if len(ns.Labels) > 0 {
		labels = "  labels:\n"
		for k, v := range ns.Labels {
			labels += fmt.Sprintf("    %s: %s\n", k, v)
		}
	}

	return fmt.Sprintf(
		"apiVersion: v1\nkind: Namespace\nmetadata:\n  name: %s\n%sstatus:\n  phase: %s\n",
		ns.Name,
		labels,
		ns.Status,
	), nil
}

func (s *NamespaceService) find(name string) (Namespace, bool) {
	for _, ns := range s.namespaces {
		if ns.Name == name {
			return ns, true
		}
	}
	return Namespace{}, false
}

func humanAge(t time.Time) string {
	hours := int(time.Since(t).Hours())
	switch {
	case hours < 1:
		return "just now"
	case hours < 24:
		return fmt.Sprintf("%dh", hours)
	default:
		return fmt.Sprintf("%dd", hours/24)
	}
}
