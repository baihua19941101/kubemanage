package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type LiveWorkloadReader struct {
	repo ClusterConnectionRepository
}

func NewLiveWorkloadReader(repo ClusterConnectionRepository) *LiveWorkloadReader {
	return &LiveWorkloadReader{repo: repo}
}

func (r *LiveWorkloadReader) ListDeployments(ctx context.Context) ([]Deployment, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()

	list, err := clientset.AppsV1().Deployments("").List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list deployments failed: %w", err)
	}
	items := make([]Deployment, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, toDeployment(item))
	}
	return items, nil
}

func (r *LiveWorkloadReader) GetDeployment(ctx context.Context, name string) (Deployment, error) {
	items, err := r.ListDeployments(ctx)
	if err != nil {
		return Deployment{}, err
	}
	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}
	return Deployment{}, fmt.Errorf("deployment not found: %s", name)
}

func (r *LiveWorkloadReader) ListPods(ctx context.Context) ([]Pod, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()

	list, err := clientset.CoreV1().Pods("").List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list pods failed: %w", err)
	}
	items := make([]Pod, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, toPod(item))
	}
	return items, nil
}

func (r *LiveWorkloadReader) ListStatefulSets(ctx context.Context) ([]StatefulSet, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()

	list, err := clientset.AppsV1().StatefulSets("").List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list statefulsets failed: %w", err)
	}
	items := make([]StatefulSet, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, toStatefulSet(item))
	}
	return items, nil
}

func (r *LiveWorkloadReader) GetStatefulSet(ctx context.Context, name string) (StatefulSet, error) {
	items, err := r.ListStatefulSets(ctx)
	if err != nil {
		return StatefulSet{}, err
	}
	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}
	return StatefulSet{}, fmt.Errorf("statefulset not found: %s", name)
}

func (r *LiveWorkloadReader) ListDaemonSets(ctx context.Context) ([]DaemonSet, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()

	list, err := clientset.AppsV1().DaemonSets("").List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list daemonsets failed: %w", err)
	}
	items := make([]DaemonSet, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, toDaemonSet(item))
	}
	return items, nil
}

func (r *LiveWorkloadReader) GetDaemonSet(ctx context.Context, name string) (DaemonSet, error) {
	items, err := r.ListDaemonSets(ctx)
	if err != nil {
		return DaemonSet{}, err
	}
	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}
	return DaemonSet{}, fmt.Errorf("daemonset not found: %s", name)
}

func (r *LiveWorkloadReader) ListJobs(ctx context.Context) ([]Job, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()

	list, err := clientset.BatchV1().Jobs("").List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list jobs failed: %w", err)
	}
	items := make([]Job, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, toJob(item))
	}
	return items, nil
}

func (r *LiveWorkloadReader) GetJob(ctx context.Context, name string) (Job, error) {
	items, err := r.ListJobs(ctx)
	if err != nil {
		return Job{}, err
	}
	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}
	return Job{}, fmt.Errorf("job not found: %s", name)
}

func (r *LiveWorkloadReader) ListCronJobs(ctx context.Context) ([]CronJob, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()

	list, err := clientset.BatchV1().CronJobs("").List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list cronjobs failed: %w", err)
	}
	items := make([]CronJob, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, toCronJob(item))
	}
	return items, nil
}

func (r *LiveWorkloadReader) GetCronJob(ctx context.Context, name string) (CronJob, error) {
	items, err := r.ListCronJobs(ctx)
	if err != nil {
		return CronJob{}, err
	}
	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}
	return CronJob{}, fmt.Errorf("cronjob not found: %s", name)
}

func (r *LiveWorkloadReader) GetPod(ctx context.Context, name string) (Pod, error) {
	items, err := r.ListPods(ctx)
	if err != nil {
		return Pod{}, err
	}
	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}
	return Pod{}, fmt.Errorf("pod not found: %s", name)
}

func (r *LiveWorkloadReader) GetPodLogs(ctx context.Context, name string, query PodLogQuery) (string, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return "", err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()

	ref, err := r.resolvePodRef(timeoutCtx, clientset, name)
	if err != nil {
		return "", err
	}

	container := strings.TrimSpace(query.Container)
	if container == "" && len(ref.containers) > 0 {
		container = ref.containers[0]
	}
	if container != "" && !containsString(ref.containers, container) {
		return "", fmt.Errorf("container not found: %s", container)
	}

	tail := int64(500)
	req := clientset.CoreV1().Pods(ref.namespace).GetLogs(ref.name, &corev1.PodLogOptions{
		Container: container,
		Follow:    false,
		TailLines: &tail,
	})
	raw, err := req.DoRaw(timeoutCtx)
	if err != nil {
		return "", fmt.Errorf("get pod logs failed: %w", err)
	}
	return applyPodLogFilter(string(raw), query), nil
}

func (r *LiveWorkloadReader) GetTerminalCapabilities(ctx context.Context, name string) (TerminalCapabilities, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return TerminalCapabilities{}, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolvePodRef(timeoutCtx, clientset, name)
	if err != nil {
		return TerminalCapabilities{}, err
	}
	return TerminalCapabilities{
		Enabled:    true,
		Protocols:  []string{"websocket"},
		Containers: ref.containers,
		Message:    "terminal bridge ready (exec endpoint pending)",
	}, nil
}

func (r *LiveWorkloadReader) CreateTerminalSession(ctx context.Context, name, container string) error {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolvePodRef(timeoutCtx, clientset, name)
	if err != nil {
		return err
	}
	target := strings.TrimSpace(container)
	if target == "" && len(ref.containers) > 0 {
		target = ref.containers[0]
	}
	if target != "" && !containsString(ref.containers, target) {
		return fmt.Errorf("container not found: %s", target)
	}
	return nil
}

func (r *LiveWorkloadReader) buildClientset(ctx context.Context) (*kubernetes.Clientset, error) {
	if r.repo == nil {
		return nil, ErrNoActiveClusterConnection
	}
	connection, err := r.repo.GetActive(ctx)
	if err != nil {
		return nil, err
	}
	cfg, err := buildRestConfig(connectionToTestInput(connection))
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("build kubernetes client failed: %w", err)
	}
	return clientset, nil
}

type podRef struct {
	namespace  string
	name       string
	containers []string
}

func (r *LiveWorkloadReader) resolvePodRef(ctx context.Context, clientset *kubernetes.Clientset, name string) (podRef, error) {
	list, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
		Limit:         20,
	})
	if err != nil {
		return podRef{}, fmt.Errorf("find pod failed: %w", err)
	}
	if len(list.Items) == 0 {
		return podRef{}, fmt.Errorf("pod not found: %s", name)
	}
	item := list.Items[0]
	containers := make([]string, 0, len(item.Spec.Containers))
	for _, c := range item.Spec.Containers {
		containers = append(containers, c.Name)
	}
	return podRef{
		namespace:  item.Namespace,
		name:       item.Name,
		containers: containers,
	}, nil
}

func containsString(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

func toDeployment(item appsv1.Deployment) Deployment {
	replicas := 0
	if item.Spec.Replicas != nil {
		replicas = int(*item.Spec.Replicas)
	}
	image := ""
	if len(item.Spec.Template.Spec.Containers) > 0 {
		image = item.Spec.Template.Spec.Containers[0].Image
	}
	return Deployment{
		Name:      item.Name,
		Namespace: item.Namespace,
		Image:     image,
		Replicas:  replicas,
		Ready:     int(item.Status.ReadyReplicas),
		CreatedAt: item.CreationTimestamp.Time,
		Age:       humanAge(item.CreationTimestamp.Time),
	}
}

func toPod(item corev1.Pod) Pod {
	image := ""
	if len(item.Spec.Containers) > 0 {
		image = item.Spec.Containers[0].Image
	}
	return Pod{
		Name:      item.Name,
		Namespace: item.Namespace,
		Node:      item.Spec.NodeName,
		Status:    string(item.Status.Phase),
		Restarts:  totalRestarts(item.Status.ContainerStatuses),
		IP:        item.Status.PodIP,
		Image:     image,
		CreatedAt: item.CreationTimestamp.Time,
		Age:       humanAge(item.CreationTimestamp.Time),
	}
}

func totalRestarts(statuses []corev1.ContainerStatus) int {
	total := 0
	for _, status := range statuses {
		total += int(status.RestartCount)
	}
	return total
}

func toStatefulSet(item appsv1.StatefulSet) StatefulSet {
	replicas := 0
	if item.Spec.Replicas != nil {
		replicas = int(*item.Spec.Replicas)
	}
	image := ""
	if len(item.Spec.Template.Spec.Containers) > 0 {
		image = item.Spec.Template.Spec.Containers[0].Image
	}
	return StatefulSet{
		Name:      item.Name,
		Namespace: item.Namespace,
		Replicas:  replicas,
		Ready:     int(item.Status.ReadyReplicas),
		Service:   item.Spec.ServiceName,
		Image:     image,
		CreatedAt: item.CreationTimestamp.Time,
		Age:       humanAge(item.CreationTimestamp.Time),
	}
}

func toDaemonSet(item appsv1.DaemonSet) DaemonSet {
	image := ""
	if len(item.Spec.Template.Spec.Containers) > 0 {
		image = item.Spec.Template.Spec.Containers[0].Image
	}
	return DaemonSet{
		Name:      item.Name,
		Namespace: item.Namespace,
		Desired:   int(item.Status.DesiredNumberScheduled),
		Current:   int(item.Status.CurrentNumberScheduled),
		Image:     image,
		CreatedAt: item.CreationTimestamp.Time,
		Age:       humanAge(item.CreationTimestamp.Time),
	}
}

func toJob(item batchv1.Job) Job {
	completions := 0
	if item.Spec.Completions != nil {
		completions = int(*item.Spec.Completions)
	}
	status := "Running"
	if item.Status.Succeeded > 0 {
		status = "Complete"
	} else if item.Status.Failed > 0 {
		status = "Failed"
	}
	return Job{
		Name:        item.Name,
		Namespace:   item.Namespace,
		Completions: completions,
		Failed:      int(item.Status.Failed),
		Status:      status,
		CreatedAt:   item.CreationTimestamp.Time,
		Age:         humanAge(item.CreationTimestamp.Time),
	}
}

func toCronJob(item batchv1.CronJob) CronJob {
	lastRun := ""
	if item.Status.LastScheduleTime != nil {
		lastRun = item.Status.LastScheduleTime.Time.Format(time.RFC3339)
	}
	return CronJob{
		Name:      item.Name,
		Namespace: item.Namespace,
		Schedule:  item.Spec.Schedule,
		Suspend:   valueOrDefault(item.Spec.Suspend, false),
		LastRun:   lastRun,
		CreatedAt: item.CreationTimestamp.Time,
		Age:       humanAge(item.CreationTimestamp.Time),
	}
}
