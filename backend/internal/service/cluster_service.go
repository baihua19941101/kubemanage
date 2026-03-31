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
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
	Nodes   int    `json:"nodes"`
}

type ClusterService struct {
	redisClient *redis.Client
	clusters    []ClusterSummary
	currentName string
}

func NewClusterService(redisClient *redis.Client) *ClusterService {
	defaultClusters := []ClusterSummary{
		{
			Name:    "demo-cluster",
			Version: "v1.30.1",
			Status:  "ready",
			Nodes:   3,
		},
		{
			Name:    "staging-cluster",
			Version: "v1.29.8",
			Status:  "ready",
			Nodes:   2,
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
