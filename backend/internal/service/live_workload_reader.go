package service

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/yaml"
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

func (r *LiveWorkloadReader) GetDeploymentYAML(ctx context.Context, name string) (string, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return "", err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolveDeploymentRef(timeoutCtx, clientset, name)
	if err != nil {
		return "", err
	}
	item, err := clientset.AppsV1().Deployments(ref.namespace).Get(timeoutCtx, ref.name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return "", fmt.Errorf("deployment not found: %s", name)
		}
		return "", fmt.Errorf("get deployment failed: %w", err)
	}
	raw, err := yaml.Marshal(item)
	if err != nil {
		return "", fmt.Errorf("marshal deployment yaml failed: %w", err)
	}
	return string(raw), nil
}

func (r *LiveWorkloadReader) UpdateDeploymentYAML(ctx context.Context, name, rawYAML string) error {
	if strings.TrimSpace(rawYAML) == "" {
		return fmt.Errorf("yaml content is empty")
	}
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolveDeploymentRef(timeoutCtx, clientset, name)
	if err != nil {
		return err
	}
	var object appsv1.Deployment
	if err := yaml.Unmarshal([]byte(rawYAML), &object); err != nil {
		return fmt.Errorf("invalid yaml: %w", err)
	}
	if strings.TrimSpace(object.Name) == "" {
		object.Name = ref.name
	}
	if strings.TrimSpace(object.Namespace) == "" {
		object.Namespace = ref.namespace
	}
	if object.Name != ref.name {
		return fmt.Errorf("deployment name mismatch: path=%s yaml=%s", ref.name, object.Name)
	}
	if object.Namespace != ref.namespace {
		return fmt.Errorf("deployment namespace mismatch: expected=%s yaml=%s", ref.namespace, object.Namespace)
	}
	if _, err := clientset.AppsV1().Deployments(ref.namespace).Update(timeoutCtx, &object, metav1.UpdateOptions{}); err != nil {
		if k8serrors.IsNotFound(err) {
			return fmt.Errorf("deployment not found: %s", ref.name)
		}
		return fmt.Errorf("update deployment failed: %w", err)
	}
	return nil
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

func (r *LiveWorkloadReader) GetStatefulSetYAML(ctx context.Context, name string) (string, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return "", err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolveStatefulSetRef(timeoutCtx, clientset, name)
	if err != nil {
		return "", err
	}
	item, err := clientset.AppsV1().StatefulSets(ref.namespace).Get(timeoutCtx, ref.name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return "", fmt.Errorf("statefulset not found: %s", name)
		}
		return "", fmt.Errorf("get statefulset failed: %w", err)
	}
	raw, err := yaml.Marshal(item)
	if err != nil {
		return "", fmt.Errorf("marshal statefulset yaml failed: %w", err)
	}
	return string(raw), nil
}

func (r *LiveWorkloadReader) UpdateStatefulSetYAML(ctx context.Context, name, rawYAML string) error {
	if strings.TrimSpace(rawYAML) == "" {
		return fmt.Errorf("yaml content is empty")
	}
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolveStatefulSetRef(timeoutCtx, clientset, name)
	if err != nil {
		return err
	}
	var object appsv1.StatefulSet
	if err := yaml.Unmarshal([]byte(rawYAML), &object); err != nil {
		return fmt.Errorf("invalid yaml: %w", err)
	}
	if strings.TrimSpace(object.Name) == "" {
		object.Name = ref.name
	}
	if strings.TrimSpace(object.Namespace) == "" {
		object.Namespace = ref.namespace
	}
	if object.Name != ref.name {
		return fmt.Errorf("statefulset name mismatch: path=%s yaml=%s", ref.name, object.Name)
	}
	if object.Namespace != ref.namespace {
		return fmt.Errorf("statefulset namespace mismatch: expected=%s yaml=%s", ref.namespace, object.Namespace)
	}
	if _, err := clientset.AppsV1().StatefulSets(ref.namespace).Update(timeoutCtx, &object, metav1.UpdateOptions{}); err != nil {
		if k8serrors.IsNotFound(err) {
			return fmt.Errorf("statefulset not found: %s", ref.name)
		}
		return fmt.Errorf("update statefulset failed: %w", err)
	}
	return nil
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

func (r *LiveWorkloadReader) GetDaemonSetYAML(ctx context.Context, name string) (string, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return "", err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolveDaemonSetRef(timeoutCtx, clientset, name)
	if err != nil {
		return "", err
	}
	item, err := clientset.AppsV1().DaemonSets(ref.namespace).Get(timeoutCtx, ref.name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return "", fmt.Errorf("daemonset not found: %s", name)
		}
		return "", fmt.Errorf("get daemonset failed: %w", err)
	}
	raw, err := yaml.Marshal(item)
	if err != nil {
		return "", fmt.Errorf("marshal daemonset yaml failed: %w", err)
	}
	return string(raw), nil
}

func (r *LiveWorkloadReader) UpdateDaemonSetYAML(ctx context.Context, name, rawYAML string) error {
	if strings.TrimSpace(rawYAML) == "" {
		return fmt.Errorf("yaml content is empty")
	}
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolveDaemonSetRef(timeoutCtx, clientset, name)
	if err != nil {
		return err
	}
	var object appsv1.DaemonSet
	if err := yaml.Unmarshal([]byte(rawYAML), &object); err != nil {
		return fmt.Errorf("invalid yaml: %w", err)
	}
	if strings.TrimSpace(object.Name) == "" {
		object.Name = ref.name
	}
	if strings.TrimSpace(object.Namespace) == "" {
		object.Namespace = ref.namespace
	}
	if object.Name != ref.name {
		return fmt.Errorf("daemonset name mismatch: path=%s yaml=%s", ref.name, object.Name)
	}
	if object.Namespace != ref.namespace {
		return fmt.Errorf("daemonset namespace mismatch: expected=%s yaml=%s", ref.namespace, object.Namespace)
	}
	if _, err := clientset.AppsV1().DaemonSets(ref.namespace).Update(timeoutCtx, &object, metav1.UpdateOptions{}); err != nil {
		if k8serrors.IsNotFound(err) {
			return fmt.Errorf("daemonset not found: %s", ref.name)
		}
		return fmt.Errorf("update daemonset failed: %w", err)
	}
	return nil
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

func (r *LiveWorkloadReader) GetJobYAML(ctx context.Context, name string) (string, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return "", err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolveJobRef(timeoutCtx, clientset, name)
	if err != nil {
		return "", err
	}
	item, err := clientset.BatchV1().Jobs(ref.namespace).Get(timeoutCtx, ref.name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return "", fmt.Errorf("job not found: %s", name)
		}
		return "", fmt.Errorf("get job failed: %w", err)
	}
	raw, err := yaml.Marshal(item)
	if err != nil {
		return "", fmt.Errorf("marshal job yaml failed: %w", err)
	}
	return string(raw), nil
}

func (r *LiveWorkloadReader) UpdateJobYAML(ctx context.Context, name, rawYAML string) error {
	if strings.TrimSpace(rawYAML) == "" {
		return fmt.Errorf("yaml content is empty")
	}
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolveJobRef(timeoutCtx, clientset, name)
	if err != nil {
		return err
	}
	var object batchv1.Job
	if err := yaml.Unmarshal([]byte(rawYAML), &object); err != nil {
		return fmt.Errorf("invalid yaml: %w", err)
	}
	if strings.TrimSpace(object.Name) == "" {
		object.Name = ref.name
	}
	if strings.TrimSpace(object.Namespace) == "" {
		object.Namespace = ref.namespace
	}
	if object.Name != ref.name {
		return fmt.Errorf("job name mismatch: path=%s yaml=%s", ref.name, object.Name)
	}
	if object.Namespace != ref.namespace {
		return fmt.Errorf("job namespace mismatch: expected=%s yaml=%s", ref.namespace, object.Namespace)
	}
	if _, err := clientset.BatchV1().Jobs(ref.namespace).Update(timeoutCtx, &object, metav1.UpdateOptions{}); err != nil {
		if k8serrors.IsNotFound(err) {
			return fmt.Errorf("job not found: %s", ref.name)
		}
		return fmt.Errorf("update job failed: %w", err)
	}
	return nil
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

func (r *LiveWorkloadReader) GetCronJobYAML(ctx context.Context, name string) (string, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return "", err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolveCronJobRef(timeoutCtx, clientset, name)
	if err != nil {
		return "", err
	}
	item, err := clientset.BatchV1().CronJobs(ref.namespace).Get(timeoutCtx, ref.name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return "", fmt.Errorf("cronjob not found: %s", name)
		}
		return "", fmt.Errorf("get cronjob failed: %w", err)
	}
	raw, err := yaml.Marshal(item)
	if err != nil {
		return "", fmt.Errorf("marshal cronjob yaml failed: %w", err)
	}
	return string(raw), nil
}

func (r *LiveWorkloadReader) UpdateCronJobYAML(ctx context.Context, name, rawYAML string) error {
	if strings.TrimSpace(rawYAML) == "" {
		return fmt.Errorf("yaml content is empty")
	}
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolveCronJobRef(timeoutCtx, clientset, name)
	if err != nil {
		return err
	}
	var object batchv1.CronJob
	if err := yaml.Unmarshal([]byte(rawYAML), &object); err != nil {
		return fmt.Errorf("invalid yaml: %w", err)
	}
	if strings.TrimSpace(object.Name) == "" {
		object.Name = ref.name
	}
	if strings.TrimSpace(object.Namespace) == "" {
		object.Namespace = ref.namespace
	}
	if object.Name != ref.name {
		return fmt.Errorf("cronjob name mismatch: path=%s yaml=%s", ref.name, object.Name)
	}
	if object.Namespace != ref.namespace {
		return fmt.Errorf("cronjob namespace mismatch: expected=%s yaml=%s", ref.namespace, object.Namespace)
	}
	if _, err := clientset.BatchV1().CronJobs(ref.namespace).Update(timeoutCtx, &object, metav1.UpdateOptions{}); err != nil {
		if k8serrors.IsNotFound(err) {
			return fmt.Errorf("cronjob not found: %s", ref.name)
		}
		return fmt.Errorf("update cronjob failed: %w", err)
	}
	return nil
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

func (r *LiveWorkloadReader) GetPodYAML(ctx context.Context, name string) (string, error) {
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
	item, err := clientset.CoreV1().Pods(ref.namespace).Get(timeoutCtx, ref.name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return "", fmt.Errorf("pod not found: %s", name)
		}
		return "", fmt.Errorf("get pod failed: %w", err)
	}
	raw, err := yaml.Marshal(item)
	if err != nil {
		return "", fmt.Errorf("marshal pod yaml failed: %w", err)
	}
	return string(raw), nil
}

func (r *LiveWorkloadReader) UpdatePodYAML(ctx context.Context, name, rawYAML string) error {
	if strings.TrimSpace(rawYAML) == "" {
		return fmt.Errorf("yaml content is empty")
	}
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
	var object corev1.Pod
	if err := yaml.Unmarshal([]byte(rawYAML), &object); err != nil {
		return fmt.Errorf("invalid yaml: %w", err)
	}
	if strings.TrimSpace(object.Name) == "" {
		object.Name = ref.name
	}
	if strings.TrimSpace(object.Namespace) == "" {
		object.Namespace = ref.namespace
	}
	if object.Name != ref.name {
		return fmt.Errorf("pod name mismatch: path=%s yaml=%s", ref.name, object.Name)
	}
	if object.Namespace != ref.namespace {
		return fmt.Errorf("pod namespace mismatch: expected=%s yaml=%s", ref.namespace, object.Namespace)
	}
	if _, err := clientset.CoreV1().Pods(ref.namespace).Update(timeoutCtx, &object, metav1.UpdateOptions{}); err != nil {
		if k8serrors.IsNotFound(err) {
			return fmt.Errorf("pod not found: %s", ref.name)
		}
		return fmt.Errorf("update pod failed: %w", err)
	}
	return nil
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

func (r *LiveWorkloadReader) ExecuteTerminal(
	ctx context.Context,
	name string,
	container string,
	command []string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	tty bool,
) error {
	if len(command) == 0 {
		command = []string{"sh"}
	}
	if r.repo == nil {
		return ErrNoActiveClusterConnection
	}
	connection, err := r.repo.GetActive(ctx)
	if err != nil {
		return err
	}
	cfg, err := buildRestConfig(connectionToTestInput(connection))
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("build kubernetes client failed: %w", err)
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

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(ref.name).
		Namespace(ref.namespace).
		SubResource("exec")
	req.VersionedParams(&corev1.PodExecOptions{
		Container: target,
		Command:   command,
		Stdin:     stdin != nil,
		Stdout:    stdout != nil,
		Stderr:    stderr != nil && !tty,
		TTY:       tty,
	}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("build terminal executor failed: %w", err)
	}
	if err := executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    tty,
	}); err != nil {
		return fmt.Errorf("exec terminal stream failed: %w", err)
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

func (r *LiveWorkloadReader) PodNamespace(ctx context.Context, name string) (string, error) {
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
	return ref.namespace, nil
}

func (r *LiveWorkloadReader) DeploymentNamespace(ctx context.Context, name string) (string, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return "", err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolveDeploymentRef(timeoutCtx, clientset, name)
	if err != nil {
		return "", err
	}
	return ref.namespace, nil
}

func (r *LiveWorkloadReader) StatefulSetNamespace(ctx context.Context, name string) (string, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return "", err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolveStatefulSetRef(timeoutCtx, clientset, name)
	if err != nil {
		return "", err
	}
	return ref.namespace, nil
}

func (r *LiveWorkloadReader) DaemonSetNamespace(ctx context.Context, name string) (string, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return "", err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolveDaemonSetRef(timeoutCtx, clientset, name)
	if err != nil {
		return "", err
	}
	return ref.namespace, nil
}

func (r *LiveWorkloadReader) JobNamespace(ctx context.Context, name string) (string, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return "", err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolveJobRef(timeoutCtx, clientset, name)
	if err != nil {
		return "", err
	}
	return ref.namespace, nil
}

func (r *LiveWorkloadReader) CronJobNamespace(ctx context.Context, name string) (string, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return "", err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ref, err := r.resolveCronJobRef(timeoutCtx, clientset, name)
	if err != nil {
		return "", err
	}
	return ref.namespace, nil
}

type podRef struct {
	namespace  string
	name       string
	containers []string
}

type deploymentRef struct {
	namespace string
	name      string
}

type statefulSetRef struct {
	namespace string
	name      string
}

type daemonSetRef struct {
	namespace string
	name      string
}

type jobRef struct {
	namespace string
	name      string
}

type cronJobRef struct {
	namespace string
	name      string
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

func (r *LiveWorkloadReader) resolveDeploymentRef(ctx context.Context, clientset *kubernetes.Clientset, name string) (deploymentRef, error) {
	list, err := clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
		Limit:         20,
	})
	if err != nil {
		return deploymentRef{}, fmt.Errorf("find deployment failed: %w", err)
	}
	if len(list.Items) == 0 {
		return deploymentRef{}, fmt.Errorf("deployment not found: %s", name)
	}
	item := list.Items[0]
	return deploymentRef{
		namespace: item.Namespace,
		name:      item.Name,
	}, nil
}

func (r *LiveWorkloadReader) resolveStatefulSetRef(ctx context.Context, clientset *kubernetes.Clientset, name string) (statefulSetRef, error) {
	list, err := clientset.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
		Limit:         20,
	})
	if err != nil {
		return statefulSetRef{}, fmt.Errorf("find statefulset failed: %w", err)
	}
	if len(list.Items) == 0 {
		return statefulSetRef{}, fmt.Errorf("statefulset not found: %s", name)
	}
	item := list.Items[0]
	return statefulSetRef{
		namespace: item.Namespace,
		name:      item.Name,
	}, nil
}

func (r *LiveWorkloadReader) resolveDaemonSetRef(ctx context.Context, clientset *kubernetes.Clientset, name string) (daemonSetRef, error) {
	list, err := clientset.AppsV1().DaemonSets("").List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
		Limit:         20,
	})
	if err != nil {
		return daemonSetRef{}, fmt.Errorf("find daemonset failed: %w", err)
	}
	if len(list.Items) == 0 {
		return daemonSetRef{}, fmt.Errorf("daemonset not found: %s", name)
	}
	item := list.Items[0]
	return daemonSetRef{
		namespace: item.Namespace,
		name:      item.Name,
	}, nil
}

func (r *LiveWorkloadReader) resolveJobRef(ctx context.Context, clientset *kubernetes.Clientset, name string) (jobRef, error) {
	list, err := clientset.BatchV1().Jobs("").List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
		Limit:         20,
	})
	if err != nil {
		return jobRef{}, fmt.Errorf("find job failed: %w", err)
	}
	if len(list.Items) == 0 {
		return jobRef{}, fmt.Errorf("job not found: %s", name)
	}
	item := list.Items[0]
	return jobRef{
		namespace: item.Namespace,
		name:      item.Name,
	}, nil
}

func (r *LiveWorkloadReader) resolveCronJobRef(ctx context.Context, clientset *kubernetes.Clientset, name string) (cronJobRef, error) {
	list, err := clientset.BatchV1().CronJobs("").List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
		Limit:         20,
	})
	if err != nil {
		return cronJobRef{}, fmt.Errorf("find cronjob failed: %w", err)
	}
	if len(list.Items) == 0 {
		return cronJobRef{}, fmt.Errorf("cronjob not found: %s", name)
	}
	item := list.Items[0]
	return cronJobRef{
		namespace: item.Namespace,
		name:      item.Name,
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
