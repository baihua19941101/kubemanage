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

type StatefulSet struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Replicas  int       `json:"replicas"`
	Ready     int       `json:"ready"`
	Service   string    `json:"service"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"-"`
	Age       string    `json:"age"`
}

type DaemonSet struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Desired   int       `json:"desired"`
	Current   int       `json:"current"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"-"`
	Age       string    `json:"age"`
}

type Job struct {
	Name        string    `json:"name"`
	Namespace   string    `json:"namespace"`
	Completions int       `json:"completions"`
	Failed      int       `json:"failed"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"-"`
	Age         string    `json:"age"`
}

type CronJob struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Schedule  string    `json:"schedule"`
	Suspend   bool      `json:"suspend"`
	LastRun   string    `json:"lastRun"`
	CreatedAt time.Time `json:"-"`
	Age       string    `json:"age"`
}

type WorkloadService struct {
	deployments    []Deployment
	pods           []Pod
	statefulSets   []StatefulSet
	daemonSets     []DaemonSet
	jobs           []Job
	cronJobs       []CronJob
	deploymentYAML map[string]string
	podYAML        map[string]string
	statefulYAML   map[string]string
	daemonYAML     map[string]string
	jobYAML        map[string]string
	cronJobYAML    map[string]string
	podLogs        map[string]string
	podLogFollow   map[string]int
}

type PodLogQuery struct {
	Keyword       string
	CaseSensitive bool
	MatchOnly     bool
	Follow        bool
}

type TerminalCapabilities struct {
	Enabled   bool     `json:"enabled"`
	Protocols []string `json:"protocols"`
	Message   string   `json:"message"`
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

	statefulSets := []StatefulSet{
		{
			Name:      "mysql",
			Namespace: "default",
			Replicas:  1,
			Ready:     1,
			Service:   "mysql-headless",
			Image:     "mysql:8.0",
			CreatedAt: now.Add(-72 * time.Hour),
		},
	}

	daemonSets := []DaemonSet{
		{
			Name:      "node-exporter",
			Namespace: "monitoring",
			Desired:   3,
			Current:   3,
			Image:     "prom/node-exporter:v1.8.0",
			CreatedAt: now.Add(-36 * time.Hour),
		},
	}

	jobs := []Job{
		{
			Name:        "db-migrate-20260401",
			Namespace:   "default",
			Completions: 1,
			Failed:      0,
			Status:      "Complete",
			CreatedAt:   now.Add(-6 * time.Hour),
		},
	}

	cronJobs := []CronJob{
		{
			Name:      "cleanup",
			Namespace: "default",
			Schedule:  "0 2 * * *",
			Suspend:   false,
			LastRun:   "2026-04-01T02:00:00+08:00",
			CreatedAt: now.Add(-240 * time.Hour),
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

	statefulYAML := map[string]string{
		"mysql": "apiVersion: apps/v1\nkind: StatefulSet\nmetadata:\n  name: mysql\n  namespace: default\nspec:\n  serviceName: mysql-headless\n  replicas: 1\n",
	}

	daemonYAML := map[string]string{
		"node-exporter": "apiVersion: apps/v1\nkind: DaemonSet\nmetadata:\n  name: node-exporter\n  namespace: monitoring\n",
	}

	jobYAML := map[string]string{
		"db-migrate-20260401": "apiVersion: batch/v1\nkind: Job\nmetadata:\n  name: db-migrate-20260401\n  namespace: default\n",
	}

	cronJobYAML := map[string]string{
		"cleanup": "apiVersion: batch/v1\nkind: CronJob\nmetadata:\n  name: cleanup\n  namespace: default\nspec:\n  schedule: \"0 2 * * *\"\n",
	}

	podLogs := map[string]string{
		"web-api-7bf59f6f9c-abcde":     "[INFO] server started on :8080\n[INFO] GET /healthz 200 1ms\n",
		"task-worker-856ddcf69f-uvwxy": "[INFO] worker tick\n[INFO] processed job id=42\n",
	}

	s := &WorkloadService{
		deployments:    deployments,
		pods:           pods,
		statefulSets:   statefulSets,
		daemonSets:     daemonSets,
		jobs:           jobs,
		cronJobs:       cronJobs,
		deploymentYAML: deployYAML,
		podYAML:        podYAML,
		statefulYAML:   statefulYAML,
		daemonYAML:     daemonYAML,
		jobYAML:        jobYAML,
		cronJobYAML:    cronJobYAML,
		podLogs:        podLogs,
		podLogFollow:   map[string]int{},
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

func (s *WorkloadService) ListStatefulSets() []StatefulSet {
	s.refreshAges()
	return slices.Clone(s.statefulSets)
}

func (s *WorkloadService) ListDaemonSets() []DaemonSet {
	s.refreshAges()
	return slices.Clone(s.daemonSets)
}

func (s *WorkloadService) ListJobs() []Job {
	s.refreshAges()
	return slices.Clone(s.jobs)
}

func (s *WorkloadService) ListCronJobs() []CronJob {
	s.refreshAges()
	return slices.Clone(s.cronJobs)
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

func (s *WorkloadService) GetStatefulSet(name string) (StatefulSet, error) {
	s.refreshAges()
	for _, item := range s.statefulSets {
		if item.Name == name {
			return item, nil
		}
	}
	return StatefulSet{}, fmt.Errorf("statefulset not found: %s", name)
}

func (s *WorkloadService) GetDaemonSet(name string) (DaemonSet, error) {
	s.refreshAges()
	for _, item := range s.daemonSets {
		if item.Name == name {
			return item, nil
		}
	}
	return DaemonSet{}, fmt.Errorf("daemonset not found: %s", name)
}

func (s *WorkloadService) GetJob(name string) (Job, error) {
	s.refreshAges()
	for _, item := range s.jobs {
		if item.Name == name {
			return item, nil
		}
	}
	return Job{}, fmt.Errorf("job not found: %s", name)
}

func (s *WorkloadService) GetCronJob(name string) (CronJob, error) {
	s.refreshAges()
	for _, item := range s.cronJobs {
		if item.Name == name {
			return item, nil
		}
	}
	return CronJob{}, fmt.Errorf("cronjob not found: %s", name)
}

func (s *WorkloadService) GetPodYAML(name string) (string, error) {
	return getYAML(s.podYAML, "pod", name)
}

func (s *WorkloadService) GetStatefulSetYAML(name string) (string, error) {
	return getYAML(s.statefulYAML, "statefulset", name)
}

func (s *WorkloadService) GetDaemonSetYAML(name string) (string, error) {
	return getYAML(s.daemonYAML, "daemonset", name)
}

func (s *WorkloadService) GetJobYAML(name string) (string, error) {
	return getYAML(s.jobYAML, "job", name)
}

func (s *WorkloadService) GetCronJobYAML(name string) (string, error) {
	return getYAML(s.cronJobYAML, "cronjob", name)
}

func (s *WorkloadService) UpdatePodYAML(name, yaml string) error {
	return updateYAML(s.podYAML, "pod", name, yaml)
}

func (s *WorkloadService) UpdateStatefulSetYAML(name, yaml string) error {
	return updateYAML(s.statefulYAML, "statefulset", name, yaml)
}

func (s *WorkloadService) UpdateDaemonSetYAML(name, yaml string) error {
	return updateYAML(s.daemonYAML, "daemonset", name, yaml)
}

func (s *WorkloadService) UpdateJobYAML(name, yaml string) error {
	return updateYAML(s.jobYAML, "job", name, yaml)
}

func (s *WorkloadService) UpdateCronJobYAML(name, yaml string) error {
	return updateYAML(s.cronJobYAML, "cronjob", name, yaml)
}

func (s *WorkloadService) GetPodLogs(name string, query PodLogQuery) (string, error) {
	logs, ok := s.podLogs[name]
	if !ok {
		return "", fmt.Errorf("pod not found: %s", name)
	}
	if query.Follow {
		s.podLogFollow[name]++
		logs = logs + fmt.Sprintf("[DEBUG] follow refresh tick=%d pod=%s\n", s.podLogFollow[name], name)
	}
	if strings.TrimSpace(query.Keyword) == "" {
		return logs, nil
	}

	lines := strings.Split(logs, "\n")
	needle := query.Keyword
	if !query.CaseSensitive {
		needle = strings.ToLower(needle)
	}

	matched := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		haystack := line
		if !query.CaseSensitive {
			haystack = strings.ToLower(line)
		}
		if strings.Contains(haystack, needle) {
			matched = append(matched, line)
		} else if !query.MatchOnly {
			matched = append(matched, line)
		}
	}

	if len(matched) == 0 {
		return "", nil
	}
	return strings.Join(matched, "\n") + "\n", nil
}

func (s *WorkloadService) GetTerminalCapabilities(name string) (TerminalCapabilities, error) {
	if _, ok := s.podLogs[name]; !ok {
		return TerminalCapabilities{}, fmt.Errorf("pod not found: %s", name)
	}
	return TerminalCapabilities{
		Enabled:   false,
		Protocols: []string{"websocket"},
		Message:   "terminal gateway not enabled",
	}, nil
}

func (s *WorkloadService) CreateTerminalSession(name string) error {
	if _, ok := s.podLogs[name]; !ok {
		return fmt.Errorf("pod not found: %s", name)
	}
	return nil
}

func (s *WorkloadService) refreshAges() {
	for i := range s.deployments {
		s.deployments[i].Age = humanAge(s.deployments[i].CreatedAt)
	}
	for i := range s.pods {
		s.pods[i].Age = humanAge(s.pods[i].CreatedAt)
	}
	for i := range s.statefulSets {
		s.statefulSets[i].Age = humanAge(s.statefulSets[i].CreatedAt)
	}
	for i := range s.daemonSets {
		s.daemonSets[i].Age = humanAge(s.daemonSets[i].CreatedAt)
	}
	for i := range s.jobs {
		s.jobs[i].Age = humanAge(s.jobs[i].CreatedAt)
	}
	for i := range s.cronJobs {
		s.cronJobs[i].Age = humanAge(s.cronJobs[i].CreatedAt)
	}
}

func getYAML(source map[string]string, kind, name string) (string, error) {
	y, ok := source[name]
	if !ok {
		return "", fmt.Errorf("%s not found: %s", kind, name)
	}
	return y, nil
}

func updateYAML(target map[string]string, kind, name, yaml string) error {
	if strings.TrimSpace(yaml) == "" {
		return fmt.Errorf("yaml content is empty")
	}
	if _, ok := target[name]; !ok {
		return fmt.Errorf("%s not found: %s", kind, name)
	}
	target[name] = yaml
	return nil
}
