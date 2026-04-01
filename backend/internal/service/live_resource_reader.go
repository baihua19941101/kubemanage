package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type LiveResourceReader struct {
	repo ClusterConnectionRepository
}

func NewLiveResourceReader(repo ClusterConnectionRepository) *LiveResourceReader {
	return &LiveResourceReader{repo: repo}
}

func (r *LiveResourceReader) ListServices(ctx context.Context) ([]ServiceItem, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	list, err := clientset.CoreV1().Services("").List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list services failed: %w", err)
	}
	items := make([]ServiceItem, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, toServiceItem(item))
	}
	return items, nil
}

func (r *LiveResourceReader) GetService(ctx context.Context, name string) (ServiceItem, error) {
	items, err := r.ListServices(ctx)
	if err != nil {
		return ServiceItem{}, err
	}
	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}
	return ServiceItem{}, fmt.Errorf("service not found: %s", name)
}

func (r *LiveResourceReader) ListConfigMaps(ctx context.Context) ([]ConfigMapItem, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	list, err := clientset.CoreV1().ConfigMaps("").List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list configmaps failed: %w", err)
	}
	items := make([]ConfigMapItem, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, ConfigMapItem{
			Name:      item.Name,
			Namespace: item.Namespace,
			DataCount: len(item.Data) + len(item.BinaryData),
			Age:       humanAge(item.CreationTimestamp.Time),
		})
	}
	return items, nil
}

func (r *LiveResourceReader) GetConfigMap(ctx context.Context, name string) (ConfigMapItem, error) {
	items, err := r.ListConfigMaps(ctx)
	if err != nil {
		return ConfigMapItem{}, err
	}
	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}
	return ConfigMapItem{}, fmt.Errorf("configmap not found: %s", name)
}

func (r *LiveResourceReader) ListSecrets(ctx context.Context) ([]SecretItem, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	list, err := clientset.CoreV1().Secrets("").List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list secrets failed: %w", err)
	}
	items := make([]SecretItem, 0, len(list.Items))
	for _, item := range list.Items {
		masked := map[string]string{}
		for key, val := range item.Data {
			masked[key] = maskSecret(string(val))
		}
		items = append(items, SecretItem{
			Name:      item.Name,
			Namespace: item.Namespace,
			Type:      string(item.Type),
			Data:      masked,
			Age:       humanAge(item.CreationTimestamp.Time),
		})
	}
	return items, nil
}

func (r *LiveResourceReader) GetSecret(ctx context.Context, name string) (SecretItem, error) {
	items, err := r.ListSecrets(ctx)
	if err != nil {
		return SecretItem{}, err
	}
	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}
	return SecretItem{}, fmt.Errorf("secret not found: %s", name)
}

func (r *LiveResourceReader) ListIngresses(ctx context.Context) ([]IngressItem, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	list, err := clientset.NetworkingV1().Ingresses("").List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list ingresses failed: %w", err)
	}
	items := make([]IngressItem, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, toIngressItem(item))
	}
	return items, nil
}

func (r *LiveResourceReader) GetIngress(ctx context.Context, name string) (IngressItem, error) {
	items, err := r.ListIngresses(ctx)
	if err != nil {
		return IngressItem{}, err
	}
	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}
	return IngressItem{}, fmt.Errorf("ingress not found: %s", name)
}

func (r *LiveResourceReader) ListIngressServices(ctx context.Context, name string) ([]ServiceItem, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	ings, err := clientset.NetworkingV1().Ingresses("").List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list ingresses failed: %w", err)
	}
	var target *networkingv1.Ingress
	for i := range ings.Items {
		if ings.Items[i].Name == name {
			target = &ings.Items[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("ingress not found: %s", name)
	}
	serviceNames := map[string]struct{}{}
	for _, rule := range target.Spec.Rules {
		if rule.HTTP == nil {
			continue
		}
		for _, path := range rule.HTTP.Paths {
			if path.Backend.Service != nil {
				serviceNames[path.Backend.Service.Name] = struct{}{}
			}
		}
	}
	if target.Spec.DefaultBackend != nil && target.Spec.DefaultBackend.Service != nil {
		serviceNames[target.Spec.DefaultBackend.Service.Name] = struct{}{}
	}
	if len(serviceNames) == 0 {
		return []ServiceItem{}, nil
	}
	services, err := clientset.CoreV1().Services(target.Namespace).List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list ingress services failed: %w", err)
	}
	items := make([]ServiceItem, 0, len(serviceNames))
	for _, svc := range services.Items {
		if _, ok := serviceNames[svc.Name]; ok {
			items = append(items, toServiceItem(svc))
		}
	}
	return items, nil
}

func (r *LiveResourceReader) ListHPAs(ctx context.Context) ([]HPAItem, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	list, err := clientset.AutoscalingV2().HorizontalPodAutoscalers("").List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list hpas failed: %w", err)
	}
	items := make([]HPAItem, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, toHPAItem(item))
	}
	return items, nil
}

func (r *LiveResourceReader) GetHPA(ctx context.Context, name string) (HPAItem, error) {
	items, err := r.ListHPAs(ctx)
	if err != nil {
		return HPAItem{}, err
	}
	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}
	return HPAItem{}, fmt.Errorf("hpa not found: %s", name)
}

func (r *LiveResourceReader) GetHPATarget(ctx context.Context, name string) (HPATarget, error) {
	item, err := r.GetHPA(ctx, name)
	if err != nil {
		return HPATarget{}, err
	}
	return HPATarget{
		Kind:            item.TargetKind,
		Name:            item.TargetName,
		Namespace:       item.Namespace,
		CurrentReplicas: item.CurrentReplicas,
		DesiredReplicas: item.CurrentReplicas,
	}, nil
}

func (r *LiveResourceReader) ListPVs(ctx context.Context) ([]PVItem, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	list, err := clientset.CoreV1().PersistentVolumes().List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list pvs failed: %w", err)
	}
	items := make([]PVItem, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, toPVItem(item))
	}
	return items, nil
}

func (r *LiveResourceReader) GetPV(ctx context.Context, name string) (PVItem, error) {
	items, err := r.ListPVs(ctx)
	if err != nil {
		return PVItem{}, err
	}
	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}
	return PVItem{}, fmt.Errorf("pv not found: %s", name)
}

func (r *LiveResourceReader) ListPVCs(ctx context.Context) ([]PVCItem, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	list, err := clientset.CoreV1().PersistentVolumeClaims("").List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list pvcs failed: %w", err)
	}
	items := make([]PVCItem, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, toPVCItem(item))
	}
	return items, nil
}

func (r *LiveResourceReader) GetPVC(ctx context.Context, name string) (PVCItem, error) {
	items, err := r.ListPVCs(ctx)
	if err != nil {
		return PVCItem{}, err
	}
	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}
	return PVCItem{}, fmt.Errorf("pvc not found: %s", name)
}

func (r *LiveResourceReader) ListStorageClasses(ctx context.Context) ([]StorageClassItem, error) {
	clientset, err := r.buildClientset(ctx)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, k8sAdapterTimeout)
	defer cancel()
	list, err := clientset.StorageV1().StorageClasses().List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list storageclasses failed: %w", err)
	}
	items := make([]StorageClassItem, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, toStorageClassItem(item))
	}
	return items, nil
}

func (r *LiveResourceReader) GetStorageClass(ctx context.Context, name string) (StorageClassItem, error) {
	items, err := r.ListStorageClasses(ctx)
	if err != nil {
		return StorageClassItem{}, err
	}
	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}
	return StorageClassItem{}, fmt.Errorf("storageclass not found: %s", name)
}

func (r *LiveResourceReader) buildClientset(ctx context.Context) (*kubernetes.Clientset, error) {
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

func toServiceItem(item corev1.Service) ServiceItem {
	ports := make([]string, 0, len(item.Spec.Ports))
	for _, p := range item.Spec.Ports {
		ports = append(ports, fmt.Sprintf("%d/%s", p.Port, string(p.Protocol)))
	}
	sort.Strings(ports)
	return ServiceItem{
		Name:      item.Name,
		Namespace: item.Namespace,
		Type:      string(item.Spec.Type),
		ClusterIP: item.Spec.ClusterIP,
		Ports:     strings.Join(ports, ","),
		Pods:      0,
		Age:       humanAge(item.CreationTimestamp.Time),
	}
}

func toIngressItem(item networkingv1.Ingress) IngressItem {
	hosts := make([]string, 0, len(item.Spec.Rules))
	for _, rule := range item.Spec.Rules {
		if strings.TrimSpace(rule.Host) != "" {
			hosts = append(hosts, rule.Host)
		}
	}
	sort.Strings(hosts)
	address := ""
	if len(item.Status.LoadBalancer.Ingress) > 0 {
		first := item.Status.LoadBalancer.Ingress[0]
		if first.IP != "" {
			address = first.IP
		} else {
			address = first.Hostname
		}
	}
	className := ""
	if item.Spec.IngressClassName != nil {
		className = *item.Spec.IngressClassName
	}
	return IngressItem{
		Name:      item.Name,
		Namespace: item.Namespace,
		ClassName: className,
		Hosts:     hosts,
		Address:   address,
		TLS:       len(item.Spec.TLS) > 0,
		Age:       humanAge(item.CreationTimestamp.Time),
	}
}

func toHPAItem(item autoscalingv2.HorizontalPodAutoscaler) HPAItem {
	targetCPU := 0
	currentCPU := 0
	for _, metric := range item.Spec.Metrics {
		if metric.Type == autoscalingv2.ResourceMetricSourceType && metric.Resource != nil && metric.Resource.Name == corev1.ResourceCPU {
			if metric.Resource.Target.AverageUtilization != nil {
				targetCPU = int(*metric.Resource.Target.AverageUtilization)
			}
		}
	}
	for _, metric := range item.Status.CurrentMetrics {
		if metric.Type == autoscalingv2.ResourceMetricSourceType && metric.Resource != nil && metric.Resource.Name == corev1.ResourceCPU {
			if metric.Resource.Current.AverageUtilization != nil {
				currentCPU = int(*metric.Resource.Current.AverageUtilization)
			}
		}
	}
	min := 0
	if item.Spec.MinReplicas != nil {
		min = int(*item.Spec.MinReplicas)
	}
	return HPAItem{
		Name:              item.Name,
		Namespace:         item.Namespace,
		TargetKind:        item.Spec.ScaleTargetRef.Kind,
		TargetName:        item.Spec.ScaleTargetRef.Name,
		MinReplicas:       min,
		MaxReplicas:       int(item.Spec.MaxReplicas),
		CurrentReplicas:   int(item.Status.CurrentReplicas),
		TargetCPUPercent:  targetCPU,
		CurrentCPUPercent: currentCPU,
		Age:               humanAge(item.CreationTimestamp.Time),
	}
}

func toPVItem(item corev1.PersistentVolume) PVItem {
	claimRef := ""
	if item.Spec.ClaimRef != nil {
		claimRef = item.Spec.ClaimRef.Namespace + "/" + item.Spec.ClaimRef.Name
	}
	capacity := ""
	if qty, ok := item.Spec.Capacity[corev1.ResourceStorage]; ok {
		capacity = qty.String()
	}
	return PVItem{
		Name:          item.Name,
		Capacity:      capacity,
		AccessModes:   joinAccessModes(item.Spec.AccessModes),
		ReclaimPolicy: string(item.Spec.PersistentVolumeReclaimPolicy),
		Status:        string(item.Status.Phase),
		ClaimRef:      claimRef,
		StorageClass:  item.Spec.StorageClassName,
		Age:           humanAge(item.CreationTimestamp.Time),
	}
}

func toPVCItem(item corev1.PersistentVolumeClaim) PVCItem {
	capacity := ""
	if qty, ok := item.Status.Capacity[corev1.ResourceStorage]; ok {
		capacity = qty.String()
	}
	return PVCItem{
		Name:         item.Name,
		Namespace:    item.Namespace,
		Status:       string(item.Status.Phase),
		Volume:       item.Spec.VolumeName,
		Capacity:     capacity,
		AccessModes:  joinAccessModes(item.Spec.AccessModes),
		StorageClass: valueOrDefault(item.Spec.StorageClassName, ""),
		Age:          humanAge(item.CreationTimestamp.Time),
	}
}

func toStorageClassItem(item storagev1.StorageClass) StorageClassItem {
	reclaim := ""
	if item.ReclaimPolicy != nil {
		reclaim = string(*item.ReclaimPolicy)
	}
	binding := ""
	if item.VolumeBindingMode != nil {
		binding = string(*item.VolumeBindingMode)
	}
	return StorageClassItem{
		Name:                 item.Name,
		Provisioner:          item.Provisioner,
		ReclaimPolicy:        reclaim,
		VolumeBindingMode:    binding,
		AllowVolumeExpansion: valueOrDefault(item.AllowVolumeExpansion, false),
		Age:                  humanAge(item.CreationTimestamp.Time),
	}
}

func joinAccessModes(modes []corev1.PersistentVolumeAccessMode) string {
	parts := make([]string, 0, len(modes))
	for _, mode := range modes {
		parts = append(parts, string(mode))
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

func valueOrDefault[T any](value *T, fallback T) T {
	if value == nil {
		return fallback
	}
	return *value
}
