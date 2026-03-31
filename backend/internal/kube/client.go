package kube

import (
	"fmt"
	"os"

	"kubeManage/backend/internal/config"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClient(cfg config.Config) (kubernetes.Interface, error) {
	if cfg.K8sMode != "live" {
		return nil, nil
	}

	restCfg, err := buildRestConfig(cfg.Kubeconfig)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, fmt.Errorf("create k8s client failed: %w", err)
	}
	return client, nil
}

func buildRestConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		if _, err := os.Stat(kubeconfig); err != nil {
			return nil, fmt.Errorf("kubeconfig not found: %s", kubeconfig)
		}
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("load kubeconfig failed: %w", err)
		}
		return cfg, nil
	}

	if env := os.Getenv("KUBECONFIG"); env != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", env)
		if err != nil {
			return nil, fmt.Errorf("load KUBECONFIG failed: %w", err)
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("load in-cluster config failed: %w", err)
	}
	return cfg, nil
}
