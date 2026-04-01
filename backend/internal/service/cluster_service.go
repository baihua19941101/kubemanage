package service

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/redis/go-redis/v9"
)

const currentClusterKey = "km:current_cluster"

type ClusterSummary struct {
	State             string `json:"state"`
	Name              string `json:"name"`
	Provider          string `json:"provider"`
	Distro            string `json:"distro"`
	KubernetesVersion string `json:"kubernetesVersion"`
	Architecture      string `json:"architecture"`
	CPU               string `json:"cpu"`
	Memory            string `json:"memory"`
	Pods              int    `json:"pods"`
}

type ClusterService struct {
	redisClient *redis.Client
	clusters    []ClusterSummary
	currentName string
}

func NewClusterService(redisClient *redis.Client) *ClusterService {
	defaultClusters := []ClusterSummary{
		{
			State:             "Ready",
			Name:              "demo-cluster",
			Provider:          "mock",
			Distro:            "mock-distro",
			KubernetesVersion: "v1.30.1",
			Architecture:      "amd64",
			CPU:               "8",
			Memory:            "16.0Gi",
			Pods:              42,
		},
		{
			State:             "Ready",
			Name:              "staging-cluster",
			Provider:          "mock",
			Distro:            "mock-distro",
			KubernetesVersion: "v1.29.8",
			Architecture:      "amd64",
			CPU:               "4",
			Memory:            "8.0Gi",
			Pods:              21,
		},
	}

	return &ClusterService{
		redisClient: redisClient,
		clusters:    defaultClusters,
		currentName: defaultClusters[0].Name,
	}
}

func (s *ClusterService) List() []ClusterSummary {
	return slices.Clone(s.clusters)
}

func (s *ClusterService) GetCurrent(ctx context.Context) (ClusterSummary, error) {
	name := s.currentName
	if s.redisClient != nil {
		redisCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		val, err := s.redisClient.Get(redisCtx, currentClusterKey).Result()
		if err == nil && val != "" {
			name = val
		}
	}

	cluster, ok := s.findByName(name)
	if ok {
		return cluster, nil
	}

	return s.clusters[0], nil
}

func (s *ClusterService) Switch(ctx context.Context, name string) error {
	if _, ok := s.findByName(name); !ok {
		return fmt.Errorf("cluster not found: %s", name)
	}

	if s.redisClient != nil {
		redisCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		if err := s.redisClient.Set(redisCtx, currentClusterKey, name, 0).Err(); err != nil {
			return fmt.Errorf("save current cluster failed: %w", err)
		}
	}

	s.currentName = name
	return nil
}

func (s *ClusterService) findByName(name string) (ClusterSummary, bool) {
	for _, c := range s.clusters {
		if c.Name == name {
			return c, true
		}
	}
	return ClusterSummary{}, false
}
