package service

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type Deployment struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Image     string    `json:"image"`
	Replicas  int       `json:"replicas"`
	Ready     int       `json:"ready"`
	CreatedAt time.Time `json:"-"`
	Age       string    `json:"age"`
}

type Pod struct {
	Name       string    `json:"name"`
	Namespace  string    `json:"namespace"`
	Node       string    `json:"node"`
	Status     string    `json:"status"`
	Restarts   int       `json:"restarts"`
	IP         string    `json:"ip"`
	Image      string    `json:"image"`
	CreatedAt  time.Time `json:"-"`
	Age        string    `json:"age"`
	deployName string
}

type WorkloadService struct {
	deployments    []Deployment
	pods           []Pod
	deploymentYAML map[string]string
	podYAML        map[string]string
	podLogs        map[string]string
}

func NewWorkloadService() *WorkloadService {
	now := time.Now()
	deployments := []Deployment{
		{
			Name:      "web-api",
			Namespace: "default",
			Image:     "nginx:1.27",
			Replicas:  3,
			Ready:     3,
			CreatedAt: now.Add(-48 * time.Hour),
		},
		{
			Name:      "task-worker",
			Namespace: "dev",
			Image:     "busybox:1.36",
			Replicas:  2,
			Ready:     2,
			CreatedAt: now.Add(-12 * time.Hour),
		},
	}

	pods := []Pod{
		{
			Name:       "web-api-7bf59f6f9c-abcde",
			Namespace:  "default",
			Node:       "node-1",
			Status:     "Running",
			Restarts:   0,
			IP:         "10.42.0.11",
			Image:      "nginx:1.27",
			CreatedAt:  now.Add(-24 * time.Hour),
			deployName: "web-api",
		},
		{
			Name:       "task-worker-856ddcf69f-uvwxy",
			Namespace:  "dev",
			Node:       "node-2",
			Status:     "Running",
			Restarts:   1,
			IP:         "10.42.1.8",
			Image:      "busybox:1.36",
			CreatedAt:  now.Add(-8 * time.Hour),
			deployName: "task-worker",
		},
	}

	deployYAML := map[string]string{
		"web-api":     "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: web-api\n  namespace: default\nspec:\n  replicas: 3\n  selector:\n    matchLabels:\n      app: web-api\n  template:\n    metadata:\n      labels:\n        app: web-api\n    spec:\n      containers:\n      - name: web-api\n        image: nginx:1.27\n",
		"task-worker": "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: task-worker\n  namespace: dev\nspec:\n  replicas: 2\n  selector:\n    matchLabels:\n      app: task-worker\n  template:\n    metadata:\n      labels:\n        app: task-worker\n    spec:\n      containers:\n      - name: task-worker\n        image: busybox:1.36\n",
	}

	podYAML := map[string]string{
		"web-api-7bf59f6f9c-abcde":     "apiVersion: v1\nkind: Pod\nmetadata:\n  name: web-api-7bf59f6f9c-abcde\n  namespace: default\nspec:\n  containers:\n  - name: web-api\n    image: nginx:1.27\n",
		"task-worker-856ddcf69f-uvwxy": "apiVersion: v1\nkind: Pod\nmetadata:\n  name: task-worker-856ddcf69f-uvwxy\n  namespace: dev\nspec:\n  containers:\n  - name: task-worker\n    image: busybox:1.36\n",
	}

	podLogs := map[string]string{
		"web-api-7bf59f6f9c-abcde":     "[INFO] server started on :8080\n[INFO] GET /healthz 200 1ms\n",
		"task-worker-856ddcf69f-uvwxy": "[INFO] worker tick\n[INFO] processed job id=42\n",
	}

	s := &WorkloadService{
		deployments:    deployments,
		pods:           pods,
		deploymentYAML: deployYAML,
		podYAML:        podYAML,
		podLogs:        podLogs,
	}
	s.refreshAges()
	return s
}

func (s *WorkloadService) ListDeployments() []Deployment {
	s.refreshAges()
	return slices.Clone(s.deployments)
}

func (s *WorkloadService) GetDeployment(name string) (Deployment, error) {
	s.refreshAges()
	for _, d := range s.deployments {
		if d.Name == name {
			return d, nil
		}
	}
	return Deployment{}, fmt.Errorf("deployment not found: %s", name)
}

func (s *WorkloadService) GetDeploymentYAML(name string) (string, error) {
	y, ok := s.deploymentYAML[name]
	if !ok {
		return "", fmt.Errorf("deployment not found: %s", name)
	}
	return y, nil
}

func (s *WorkloadService) UpdateDeploymentYAML(name, yaml string) error {
	if strings.TrimSpace(yaml) == "" {
		return fmt.Errorf("yaml content is empty")
	}
	if _, ok := s.deploymentYAML[name]; !ok {
		return fmt.Errorf("deployment not found: %s", name)
	}
	s.deploymentYAML[name] = yaml
	return nil
}

func (s *WorkloadService) ListPods() []Pod {
	s.refreshAges()
	return slices.Clone(s.pods)
}

func (s *WorkloadService) GetPod(name string) (Pod, error) {
	s.refreshAges()
	for _, p := range s.pods {
		if p.Name == name {
			return p, nil
		}
	}
	return Pod{}, fmt.Errorf("pod not found: %s", name)
}

func (s *WorkloadService) GetPodYAML(name string) (string, error) {
	y, ok := s.podYAML[name]
	if !ok {
		return "", fmt.Errorf("pod not found: %s", name)
	}
	return y, nil
}

func (s *WorkloadService) UpdatePodYAML(name, yaml string) error {
	if strings.TrimSpace(yaml) == "" {
		return fmt.Errorf("yaml content is empty")
	}
	if _, ok := s.podYAML[name]; !ok {
		return fmt.Errorf("pod not found: %s", name)
	}
	s.podYAML[name] = yaml
	return nil
}

func (s *WorkloadService) GetPodLogs(name string) (string, error) {
	logs, ok := s.podLogs[name]
	if !ok {
		return "", fmt.Errorf("pod not found: %s", name)
	}
	return logs, nil
}

func (s *WorkloadService) refreshAges() {
	for i := range s.deployments {
		s.deployments[i].Age = humanAge(s.deployments[i].CreatedAt)
	}
	for i := range s.pods {
		s.pods[i].Age = humanAge(s.pods[i].CreatedAt)
	}
}
