package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"k8s.io/client-go/tools/clientcmd"
	"kubeManage/backend/internal/infra"
)

type ClusterConnectionMode string

const (
	ClusterConnectionModeKubeconfig ClusterConnectionMode = "kubeconfig"
	ClusterConnectionModeToken      ClusterConnectionMode = "token"
)

var ErrNoActiveClusterConnection = errors.New("no active cluster connection")

type ClusterConnection struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	Mode              string    `json:"mode"`
	APIServer         string    `json:"apiServer"`
	SkipTLSVerify     bool      `json:"skipTLSVerify"`
	IsDefault         bool      `json:"isDefault"`
	Status            string    `json:"status"`
	LastCheckedAt     string    `json:"lastCheckedAt,omitempty"`
	LastError         string    `json:"lastError,omitempty"`
	HasKubeconfig     bool      `json:"hasKubeconfig"`
	HasBearerToken    bool      `json:"hasBearerToken"`
	HasCACert         bool      `json:"hasCaCert"`
	KubeconfigPreview string    `json:"kubeconfigPreview,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type ImportKubeconfigInput struct {
	Name              string
	KubeconfigContent string
}

type ImportTokenInput struct {
	Name          string
	APIServer     string
	BearerToken   string
	CACert        string
	SkipTLSVerify bool
}

type ConnectionTestInput struct {
	Mode              string
	APIServer         string
	KubeconfigContent string
	BearerToken       string
	CACert            string
	SkipTLSVerify     bool
}

type ConnectionTestResult struct {
	Success        bool   `json:"success"`
	Version        string `json:"version,omitempty"`
	Server         string `json:"server,omitempty"`
	NodeCount      int    `json:"nodeCount,omitempty"`
	NamespaceCount int    `json:"namespaceCount,omitempty"`
	Message        string `json:"message"`
}

type LiveClusterSummary struct {
	State             string `json:"state"`
	Name              string `json:"name"`
	Provider          string `json:"provider"`
	Distro            string `json:"distro"`
	KubernetesVersion string `json:"kubernetesVersion"`
	Architecture      string `json:"architecture"`
	CPU               string `json:"cpu"`
	Memory            string `json:"memory"`
	Pods              int    `json:"pods"`
	APIServer         string `json:"apiServer"`
	Source            string `json:"source"`
}

type ClusterConnectionRepository interface {
	List(ctx context.Context) ([]infra.ClusterConnectionRecord, error)
	Create(ctx context.Context, record *infra.ClusterConnectionRecord) error
	Get(ctx context.Context, id uint) (infra.ClusterConnectionRecord, error)
	SetActive(ctx context.Context, id uint) error
	GetActive(ctx context.Context) (infra.ClusterConnectionRecord, error)
	UpdateStatus(ctx context.Context, id uint, status string, checkedAt time.Time, lastError string) error
}

type K8sAdapter interface {
	TestConnection(ctx context.Context, input ConnectionTestInput) (ConnectionTestResult, error)
	GetClusterSummary(ctx context.Context, connection infra.ClusterConnectionRecord) (LiveClusterSummary, error)
	ListNamespaces(ctx context.Context, connection infra.ClusterConnectionRecord) ([]Namespace, error)
}

type ClusterConnectionService struct {
	repo    ClusterConnectionRepository
	adapter K8sAdapter
}

func NewClusterConnectionService(repo ClusterConnectionRepository, adapter K8sAdapter) *ClusterConnectionService {
	if repo == nil {
		repo = newMemoryClusterConnectionRepo()
	}
	if adapter == nil {
		adapter = &fakeK8sAdapter{}
	}
	return &ClusterConnectionService{repo: repo, adapter: adapter}
}

func (s *ClusterConnectionService) List(ctx context.Context) ([]ClusterConnection, error) {
	records, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]ClusterConnection, 0, len(records))
	for _, item := range records {
		items = append(items, sanitizeConnection(item))
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items, nil
}

func (s *ClusterConnectionService) ImportKubeconfig(ctx context.Context, input ImportKubeconfigInput) (ClusterConnection, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" || strings.TrimSpace(input.KubeconfigContent) == "" {
		return ClusterConnection{}, errors.New("name and kubeconfig content are required")
	}
	cfg, err := clientcmd.Load([]byte(input.KubeconfigContent))
	if err != nil {
		return ClusterConnection{}, fmt.Errorf("parse kubeconfig failed: %w", err)
	}
	apiServer := ""
	if cfg.CurrentContext != "" {
		if ctxCfg, ok := cfg.Contexts[cfg.CurrentContext]; ok {
			if cluster, ok := cfg.Clusters[ctxCfg.Cluster]; ok {
				apiServer = cluster.Server
			}
		}
	}
	record := infra.ClusterConnectionRecord{
		Name:              name,
		Mode:              string(ClusterConnectionModeKubeconfig),
		APIServer:         apiServer,
		KubeconfigContent: input.KubeconfigContent,
		Status:            "unknown",
	}
	if err := s.repo.Create(ctx, &record); err != nil {
		return ClusterConnection{}, err
	}
	return sanitizeConnection(record), nil
}

func (s *ClusterConnectionService) ImportToken(ctx context.Context, input ImportTokenInput) (ClusterConnection, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" || strings.TrimSpace(input.APIServer) == "" || strings.TrimSpace(input.BearerToken) == "" {
		return ClusterConnection{}, errors.New("name, api server and bearer token are required")
	}
	record := infra.ClusterConnectionRecord{
		Name:          name,
		Mode:          string(ClusterConnectionModeToken),
		APIServer:     strings.TrimSpace(input.APIServer),
		BearerToken:   input.BearerToken,
		CACert:        input.CACert,
		SkipTLSVerify: input.SkipTLSVerify,
		Status:        "unknown",
	}
	if err := s.repo.Create(ctx, &record); err != nil {
		return ClusterConnection{}, err
	}
	return sanitizeConnection(record), nil
}

func (s *ClusterConnectionService) TestConnection(ctx context.Context, input ConnectionTestInput) (ConnectionTestResult, error) {
	return s.adapter.TestConnection(ctx, input)
}

func (s *ClusterConnectionService) Activate(ctx context.Context, id uint) error {
	record, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	result, testErr := s.adapter.TestConnection(ctx, ConnectionTestInput{
		Mode:              record.Mode,
		APIServer:         record.APIServer,
		KubeconfigContent: record.KubeconfigContent,
		BearerToken:       record.BearerToken,
		CACert:            record.CACert,
		SkipTLSVerify:     record.SkipTLSVerify,
	})
	checkedAt := time.Now()
	if testErr != nil {
		_ = s.repo.UpdateStatus(ctx, id, "failed", checkedAt, testErr.Error())
		return testErr
	}
	status := "connected"
	if !result.Success {
		status = "failed"
	}
	if err := s.repo.UpdateStatus(ctx, id, status, checkedAt, result.Message); err != nil {
		return err
	}
	if !result.Success {
		return errors.New(result.Message)
	}
	return s.repo.SetActive(ctx, id)
}

func (s *ClusterConnectionService) GetLiveCluster(ctx context.Context) (LiveClusterSummary, error) {
	record, err := s.repo.GetActive(ctx)
	if err != nil {
		return LiveClusterSummary{}, err
	}
	return s.adapter.GetClusterSummary(ctx, record)
}

func (s *ClusterConnectionService) ListLiveNamespaces(ctx context.Context) ([]Namespace, error) {
	record, err := s.repo.GetActive(ctx)
	if err != nil {
		return nil, err
	}
	return s.adapter.ListNamespaces(ctx, record)
}

func sanitizeConnection(item infra.ClusterConnectionRecord) ClusterConnection {
	out := ClusterConnection{
		ID:             item.ID,
		Name:           item.Name,
		Mode:           item.Mode,
		APIServer:      item.APIServer,
		SkipTLSVerify:  item.SkipTLSVerify,
		IsDefault:      item.IsDefault,
		Status:         item.Status,
		LastError:      item.LastError,
		HasKubeconfig:  strings.TrimSpace(item.KubeconfigContent) != "",
		HasBearerToken: strings.TrimSpace(item.BearerToken) != "",
		HasCACert:      strings.TrimSpace(item.CACert) != "",
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
	if item.LastCheckedAt != nil {
		out.LastCheckedAt = item.LastCheckedAt.Format(time.RFC3339)
	}
	if out.HasKubeconfig {
		out.KubeconfigPreview = "***"
	}
	return out
}

type memoryClusterConnectionRepo struct {
	mu      sync.RWMutex
	nextID  uint
	records []infra.ClusterConnectionRecord
}

func newMemoryClusterConnectionRepo() *memoryClusterConnectionRepo {
	return &memoryClusterConnectionRepo{nextID: 1, records: make([]infra.ClusterConnectionRecord, 0, 4)}
}

func (r *memoryClusterConnectionRepo) List(_ context.Context) ([]infra.ClusterConnectionRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]infra.ClusterConnectionRecord, len(r.records))
	copy(out, r.records)
	return out, nil
}

func (r *memoryClusterConnectionRepo) Create(_ context.Context, record *infra.ClusterConnectionRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, item := range r.records {
		if item.Name == record.Name {
			return fmt.Errorf("cluster connection already exists: %s", record.Name)
		}
	}
	now := time.Now()
	record.ID = r.nextID
	record.CreatedAt = now
	record.UpdatedAt = now
	r.nextID++
	r.records = append(r.records, *record)
	return nil
}

func (r *memoryClusterConnectionRepo) Get(_ context.Context, id uint) (infra.ClusterConnectionRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, item := range r.records {
		if item.ID == id {
			return item, nil
		}
	}
	return infra.ClusterConnectionRecord{}, fmt.Errorf("cluster connection not found: %d", id)
}

func (r *memoryClusterConnectionRepo) SetActive(_ context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	found := false
	for i := range r.records {
		r.records[i].IsDefault = r.records[i].ID == id
		if r.records[i].ID == id {
			r.records[i].UpdatedAt = time.Now()
			found = true
		}
	}
	if !found {
		return fmt.Errorf("cluster connection not found: %d", id)
	}
	return nil
}

func (r *memoryClusterConnectionRepo) GetActive(_ context.Context) (infra.ClusterConnectionRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, item := range r.records {
		if item.IsDefault {
			return item, nil
		}
	}
	return infra.ClusterConnectionRecord{}, ErrNoActiveClusterConnection
}

func (r *memoryClusterConnectionRepo) UpdateStatus(_ context.Context, id uint, status string, checkedAt time.Time, lastError string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.records {
		if r.records[i].ID == id {
			r.records[i].Status = status
			r.records[i].LastCheckedAt = &checkedAt
			r.records[i].LastError = lastError
			r.records[i].UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("cluster connection not found: %d", id)
}
