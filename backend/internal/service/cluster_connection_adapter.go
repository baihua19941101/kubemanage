package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"kubeManage/backend/internal/infra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const k8sAdapterTimeout = 8 * time.Second

type fakeK8sAdapter struct{}

func (a *fakeK8sAdapter) TestConnection(_ context.Context, input ConnectionTestInput) (ConnectionTestResult, error) {
	if strings.TrimSpace(input.Mode) == "" {
		return ConnectionTestResult{}, errors.New("connection mode is required")
	}
	server := input.APIServer
	if strings.TrimSpace(input.KubeconfigContent) != "" {
		server = "kubeconfig-server"
	}
	return ConnectionTestResult{
		Success:        true,
		Version:        "v1.30.0",
		Server:         server,
		NodeCount:      1,
		NamespaceCount: 2,
		Message:        "connection ok",
	}, nil
}

func (a *fakeK8sAdapter) GetClusterSummary(_ context.Context, connection infra.ClusterConnectionRecord) (LiveClusterSummary, error) {
	return LiveClusterSummary{
		Name:      connection.Name,
		Version:   "v1.30.0",
		Status:    "ready",
		Nodes:     1,
		APIServer: connection.APIServer,
		Source:    "fake",
	}, nil
}

func (a *fakeK8sAdapter) ListNamespaces(_ context.Context, _ infra.ClusterConnectionRecord) ([]Namespace, error) {
	now := time.Now()
	return []Namespace{
		{Name: "default", Status: "Active", CreatedAt: now.Add(-48 * time.Hour), Age: humanAge(now.Add(-48 * time.Hour))},
		{Name: "kube-system", Status: "Active", CreatedAt: now.Add(-72 * time.Hour), Age: humanAge(now.Add(-72 * time.Hour))},
	}, nil
}

type realK8sAdapter struct{}

func NewRealK8sAdapter() K8sAdapter { return &realK8sAdapter{} }

func (a *realK8sAdapter) TestConnection(ctx context.Context, input ConnectionTestInput) (ConnectionTestResult, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()

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
	nodes, err := clientset.CoreV1().Nodes().List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return ConnectionTestResult{}, fmt.Errorf("list nodes failed: %w", err)
	}
	namespaces, err := clientset.CoreV1().Namespaces().List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return ConnectionTestResult{}, fmt.Errorf("list namespaces failed: %w", err)
	}
	return ConnectionTestResult{
		Success:        true,
		Version:        version.GitVersion,
		Server:         cfg.Host,
		NodeCount:      len(nodes.Items),
		NamespaceCount: len(namespaces.Items),
		Message:        "connection ok",
	}, nil
}

func (a *realK8sAdapter) GetClusterSummary(ctx context.Context, connection infra.ClusterConnectionRecord) (LiveClusterSummary, error) {
	result, err := a.TestConnection(ctx, connectionToTestInput(connection))
	if err != nil {
		return LiveClusterSummary{}, err
	}
	return LiveClusterSummary{
		Name:      connection.Name,
		Version:   result.Version,
		Status:    "ready",
		Nodes:     result.NodeCount,
		APIServer: result.Server,
		Source:    "live",
	}, nil
}

func (a *realK8sAdapter) ListNamespaces(ctx context.Context, connection infra.ClusterConnectionRecord) ([]Namespace, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()

	cfg, err := buildRestConfig(connectionToTestInput(connection))
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("build kubernetes client failed: %w", err)
	}
	list, err := clientset.CoreV1().Namespaces().List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list namespaces failed: %w", err)
	}
	items := make([]Namespace, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, Namespace{
			Name:      item.Name,
			Status:    string(item.Status.Phase),
			Labels:    item.Labels,
			CreatedAt: item.CreationTimestamp.Time,
			Age:       humanAge(item.CreationTimestamp.Time),
		})
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
		return &rest.Config{
			Host:        strings.TrimSpace(input.APIServer),
			BearerToken: input.BearerToken,
			TLSClientConfig: rest.TLSClientConfig{
				CAData:   []byte(input.CACert),
				Insecure: input.SkipTLSVerify,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported connection mode: %s", input.Mode)
	}
}

func connectionToTestInput(connection infra.ClusterConnectionRecord) ConnectionTestInput {
	return ConnectionTestInput{
		Mode:              connection.Mode,
		APIServer:         connection.APIServer,
		KubeconfigContent: connection.KubeconfigContent,
		BearerToken:       connection.BearerToken,
		CACert:            connection.CACert,
		SkipTLSVerify:     connection.SkipTLSVerify,
	}
}
