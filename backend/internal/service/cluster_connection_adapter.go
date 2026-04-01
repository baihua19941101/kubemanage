package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"kubeManage/backend/internal/infra"

	corev1 "k8s.io/api/core/v1"
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
		State:             "Ready",
		Name:              connection.Name,
		Provider:          "mock",
		Distro:            "mock-distro",
		KubernetesVersion: "v1.30.0",
		Architecture:      "amd64",
		CPU:               "4",
		Memory:            "8.0Gi",
		Pods:              12,
		APIServer:         connection.APIServer,
		Source:            "fake",
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
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()

	cfg, err := buildRestConfig(connectionToTestInput(connection))
	if err != nil {
		return LiveClusterSummary{}, err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return LiveClusterSummary{}, fmt.Errorf("build kubernetes client failed: %w", err)
	}
	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return LiveClusterSummary{}, fmt.Errorf("get server version failed: %w", err)
	}
	nodes, err := clientset.CoreV1().Nodes().List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return LiveClusterSummary{}, fmt.Errorf("list nodes failed: %w", err)
	}
	pods, err := clientset.CoreV1().Pods("").List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return LiveClusterSummary{}, fmt.Errorf("list pods failed: %w", err)
	}

	return LiveClusterSummary{
		State:             clusterState(nodes.Items),
		Name:              connection.Name,
		Provider:          detectProvider(nodes.Items),
		Distro:            detectDistro(nodes.Items),
		KubernetesVersion: version.GitVersion,
		Architecture:      detectArchitecture(nodes.Items),
		CPU:               totalAllocatableCPU(nodes.Items),
		Memory:            totalAllocatableMemory(nodes.Items),
		Pods:              len(pods.Items),
		APIServer:         cfg.Host,
		Source:            "live",
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

func clusterState(nodes []corev1.Node) string {
	if len(nodes) == 0 {
		return "Unknown"
	}
	ready := 0
	for _, node := range nodes {
		for _, cond := range node.Status.Conditions {
			if cond.Type == corev1.NodeReady && cond.Status == corev1.ConditionTrue {
				ready++
				break
			}
		}
	}
	switch {
	case ready == len(nodes):
		return "Ready"
	case ready == 0:
		return "NotReady"
	default:
		return "Degraded"
	}
}

func detectProvider(nodes []corev1.Node) string {
	if len(nodes) == 0 {
		return "Unknown"
	}
	seen := map[string]struct{}{}
	for _, node := range nodes {
		providerID := strings.ToLower(strings.TrimSpace(node.Spec.ProviderID))
		if providerID == "" {
			seen["OnPrem"] = struct{}{}
			continue
		}
		prefix := providerID
		if idx := strings.Index(prefix, "://"); idx >= 0 {
			prefix = prefix[:idx]
		}
		switch prefix {
		case "aws":
			seen["AWS"] = struct{}{}
		case "gce", "gke":
			seen["GCP"] = struct{}{}
		case "azure":
			seen["Azure"] = struct{}{}
		case "aliyun":
			seen["Aliyun"] = struct{}{}
		case "vsphere":
			seen["vSphere"] = struct{}{}
		case "openstack":
			seen["OpenStack"] = struct{}{}
		default:
			seen[strings.ToUpper(prefix)] = struct{}{}
		}
	}
	return collapseSingleOrMixed(seen)
}

func detectDistro(nodes []corev1.Node) string {
	if len(nodes) == 0 {
		return "Unknown"
	}
	seen := map[string]struct{}{}
	for _, node := range nodes {
		osImage := strings.TrimSpace(node.Status.NodeInfo.OSImage)
		if osImage == "" {
			continue
		}
		seen[osImage] = struct{}{}
	}
	if len(seen) == 0 {
		return "Unknown"
	}
	return collapseSingleOrMixed(seen)
}

func detectArchitecture(nodes []corev1.Node) string {
	if len(nodes) == 0 {
		return "Unknown"
	}
	seen := map[string]struct{}{}
	for _, node := range nodes {
		arch := strings.TrimSpace(node.Status.NodeInfo.Architecture)
		if arch == "" {
			continue
		}
		seen[arch] = struct{}{}
	}
	if len(seen) == 0 {
		return "Unknown"
	}
	return collapseSingleOrMixed(seen)
}

func totalAllocatableCPU(nodes []corev1.Node) string {
	var totalMilli int64
	for _, node := range nodes {
		totalMilli += node.Status.Allocatable.Cpu().MilliValue()
	}
	if totalMilli%1000 == 0 {
		return strconv.FormatInt(totalMilli/1000, 10)
	}
	return fmt.Sprintf("%.1f", float64(totalMilli)/1000.0)
}

func totalAllocatableMemory(nodes []corev1.Node) string {
	var totalBytes int64
	for _, node := range nodes {
		totalBytes += node.Status.Allocatable.Memory().Value()
	}
	return fmt.Sprintf("%.1fGi", float64(totalBytes)/(1024*1024*1024))
}

func collapseSingleOrMixed(values map[string]struct{}) string {
	if len(values) == 0 {
		return "Unknown"
	}
	parts := make([]string, 0, len(values))
	for v := range values {
		parts = append(parts, v)
	}
	sort.Strings(parts)
	if len(parts) == 1 {
		return parts[0]
	}
	return "Mixed"
}
