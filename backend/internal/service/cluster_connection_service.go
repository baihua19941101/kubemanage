package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"kubeManage/backend/internal/infra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterConnectionMode string

const (
	ClusterConnectionModeKubeconfig ClusterConnectionMode = "kubeconfig"
	ClusterConnectionModeToken      ClusterConnectionMode = "token"
)

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
	Name      string `json:"name"`
	Version   string `json:"version"`
	Status    string `json:"status"`
	Nodes     int    `json:"nodes"`
	APIServer string `json:"apiServer"`
	Source    string `json:"source"`
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
	return infra.ClusterConnectionRecord{}, errors.New("no active cluster connection")
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

type fakeK8sAdapter struct{}

func (a *fakeK8sAdapter) TestConnection(_ context.Context, input ConnectionTestInput) (ConnectionTestResult, error) {
	if strings.TrimSpace(input.Mode) == "" {
		return ConnectionTestResult{}, errors.New("connection mode is required")
	}
	server := input.APIServer
	if strings.TrimSpace(input.KubeconfigContent) != "" {
		server = "kubeconfig-server"
	}
	return ConnectionTestResult{Success: true, Version: "v1.30.0", Server: server, NodeCount: 1, NamespaceCount: 2, Message: "connection ok"}, nil
}

func (a *fakeK8sAdapter) GetClusterSummary(_ context.Context, connection infra.ClusterConnectionRecord) (LiveClusterSummary, error) {
	return LiveClusterSummary{Name: connection.Name, Version: "v1.30.0", Status: "ready", Nodes: 1, APIServer: connection.APIServer, Source: "fake"}, nil
}

func (a *fakeK8sAdapter) ListNamespaces(_ context.Context, _ infra.ClusterConnectionRecord) ([]Namespace, error) {
	now := time.Now()
	return []Namespace{{Name: "default", Status: "Active", CreatedAt: now.Add(-48 * time.Hour), Age: humanAge(now.Add(-48 * time.Hour))}, {Name: "kube-system", Status: "Active", CreatedAt: now.Add(-72 * time.Hour), Age: humanAge(now.Add(-72 * time.Hour))}}, nil
}

type realK8sAdapter struct{}

func NewRealK8sAdapter() K8sAdapter { return &realK8sAdapter{} }

func (a *realK8sAdapter) TestConnection(ctx context.Context, input ConnectionTestInput) (ConnectionTestResult, error) {
	cfg, err := buildRestConfig(input)
	if err != nil {
		return ConnectionTestResult{}, err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return ConnectionTestResult{}, fmt.Errorf("build kubernetes client failed: %w", err)
	}
	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return ConnectionTestResult{}, fmt.Errorf("get server version failed: %w", err)
	}
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return ConnectionTestResult{}, fmt.Errorf("list nodes failed: %w", err)
	}
	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return ConnectionTestResult{}, fmt.Errorf("list namespaces failed: %w", err)
	}
	return ConnectionTestResult{Success: true, Version: version.GitVersion, Server: cfg.Host, NodeCount: len(nodes.Items), NamespaceCount: len(namespaces.Items), Message: "connection ok"}, nil
}

func (a *realK8sAdapter) GetClusterSummary(ctx context.Context, connection infra.ClusterConnectionRecord) (LiveClusterSummary, error) {
	result, err := a.TestConnection(ctx, ConnectionTestInput{Mode: connection.Mode, APIServer: connection.APIServer, KubeconfigContent: connection.KubeconfigContent, BearerToken: connection.BearerToken, CACert: connection.CACert, SkipTLSVerify: connection.SkipTLSVerify})
	if err != nil {
		return LiveClusterSummary{}, err
	}
	return LiveClusterSummary{Name: connection.Name, Version: result.Version, Status: "ready", Nodes: result.NodeCount, APIServer: result.Server, Source: "live"}, nil
}

func (a *realK8sAdapter) ListNamespaces(ctx context.Context, connection infra.ClusterConnectionRecord) ([]Namespace, error) {
	cfg, err := buildRestConfig(ConnectionTestInput{Mode: connection.Mode, APIServer: connection.APIServer, KubeconfigContent: connection.KubeconfigContent, BearerToken: connection.BearerToken, CACert: connection.CACert, SkipTLSVerify: connection.SkipTLSVerify})
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("build kubernetes client failed: %w", err)
	}
	list, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list namespaces failed: %w", err)
	}
	items := make([]Namespace, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, Namespace{Name: item.Name, Status: string(item.Status.Phase), Labels: item.Labels, CreatedAt: item.CreationTimestamp.Time, Age: humanAge(item.CreationTimestamp.Time)})
	}
	return items, nil
}

func buildRestConfig(input ConnectionTestInput) (*rest.Config, error) {
	switch input.Mode {
	case string(ClusterConnectionModeKubeconfig):
		if strings.TrimSpace(input.KubeconfigContent) == "" {
			return nil, errors.New("kubeconfig content is required")
		}
		cfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(input.KubeconfigContent))
		if err != nil {
			return nil, fmt.Errorf("build kubeconfig rest config failed: %w", err)
		}
		return cfg, nil
	case string(ClusterConnectionModeToken):
		if strings.TrimSpace(input.APIServer) == "" || strings.TrimSpace(input.BearerToken) == "" {
			return nil, errors.New("api server and bearer token are required")
		}
		return &rest.Config{Host: strings.TrimSpace(input.APIServer), BearerToken: input.BearerToken, TLSClientConfig: rest.TLSClientConfig{CAData: []byte(input.CACert), Insecure: input.SkipTLSVerify}}, nil
	default:
		return nil, fmt.Errorf("unsupported connection mode: %s", input.Mode)
	}
}
