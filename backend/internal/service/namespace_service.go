package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
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
	kubeClient kubernetes.Interface
	liveMode   bool
}

func NewNamespaceService(kubeClient kubernetes.Interface, liveMode bool) *NamespaceService {
	now := time.Now()
	return &NamespaceService{
		namespaces: []Namespace{
			{Name: "default", Status: "Active", CreatedAt: now.Add(-240 * time.Hour), Labels: map[string]string{"system": "true"}},
			{Name: "kube-system", Status: "Active", CreatedAt: now.Add(-240 * time.Hour), Labels: map[string]string{"system": "true"}},
			{Name: "dev", Status: "Active", CreatedAt: now.Add(-24 * time.Hour), Labels: map[string]string{"env": "dev"}},
		},
		kubeClient: kubeClient,
		liveMode:   liveMode,
	}
}

func (s *NamespaceService) List() []Namespace {
	if s.liveMode && s.kubeClient != nil {
		items, err := s.listFromK8s(context.Background())
		if err == nil {
			return items
		}
	}
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

	if s.liveMode && s.kubeClient != nil {
		return s.createInK8s(context.Background(), name, labels)
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
	if s.liveMode && s.kubeClient != nil {
		return s.kubeClient.CoreV1().Namespaces().Delete(context.Background(), name, metav1.DeleteOptions{})
	}

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
	if s.liveMode && s.kubeClient != nil {
		ns, err := s.kubeClient.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return "", fmt.Errorf("namespace not found: %s", name)
		}
		b, err := json.Marshal(ns)
		if err != nil {
			return "", fmt.Errorf("marshal namespace failed: %w", err)
		}
		y, err := yaml.JSONToYAML(b)
		if err != nil {
			return "", fmt.Errorf("yaml encode failed: %w", err)
		}
		return string(y), nil
	}

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

func (s *NamespaceService) listFromK8s(ctx context.Context) ([]Namespace, error) {
	list, err := s.kubeClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]Namespace, 0, len(list.Items))
	for _, ns := range list.Items {
		items = append(items, Namespace{
			Name:      ns.Name,
			Status:    string(ns.Status.Phase),
			Age:       humanAge(ns.CreationTimestamp.Time),
			Labels:    ns.Labels,
			CreatedAt: ns.CreationTimestamp.Time,
		})
	}
	return items, nil
}

func (s *NamespaceService) createInK8s(ctx context.Context, name string, labels map[string]string) (Namespace, error) {
	if _, err := s.kubeClient.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{}); err == nil {
		return Namespace{}, fmt.Errorf("namespace already exists: %s", name)
	}

	created, err := s.kubeClient.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return Namespace{}, err
	}
	return Namespace{
		Name:      created.Name,
		Status:    string(created.Status.Phase),
		Age:       humanAge(created.CreationTimestamp.Time),
		Labels:    created.Labels,
		CreatedAt: created.CreationTimestamp.Time,
	}, nil
}
