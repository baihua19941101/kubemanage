package service

import (
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"
)

type ServiceItem struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	ClusterIP string `json:"clusterIP"`
	Ports     string `json:"ports"`
	Pods      int    `json:"pods"`
	Age       string `json:"age"`
}

type ConfigMapItem struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	DataCount int    `json:"dataCount"`
	Age       string `json:"age"`
}

type SecretItem struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Type      string            `json:"type"`
	Data      map[string]string `json:"data"`
	Age       string            `json:"age"`
}

type IngressItem struct {
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	ClassName string   `json:"className"`
	Hosts     []string `json:"hosts"`
	Address   string   `json:"address"`
	TLS       bool     `json:"tls"`
	Age       string   `json:"age"`
}

type HPAItem struct {
	Name              string `json:"name"`
	Namespace         string `json:"namespace"`
	TargetKind        string `json:"targetKind"`
	TargetName        string `json:"targetName"`
	MinReplicas       int    `json:"minReplicas"`
	MaxReplicas       int    `json:"maxReplicas"`
	CurrentReplicas   int    `json:"currentReplicas"`
	TargetCPUPercent  int    `json:"targetCPUPercent"`
	CurrentCPUPercent int    `json:"currentCPUPercent"`
	Age               string `json:"age"`
}

type HPATarget struct {
	Kind            string `json:"kind"`
	Name            string `json:"name"`
	Namespace       string `json:"namespace"`
	CurrentReplicas int    `json:"currentReplicas"`
	DesiredReplicas int    `json:"desiredReplicas"`
}

type PVItem struct {
	Name          string `json:"name"`
	Capacity      string `json:"capacity"`
	AccessModes   string `json:"accessModes"`
	ReclaimPolicy string `json:"reclaimPolicy"`
	Status        string `json:"status"`
	ClaimRef      string `json:"claimRef"`
	StorageClass  string `json:"storageClass"`
	Age           string `json:"age"`
}

type PVCItem struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	Status       string `json:"status"`
	Volume       string `json:"volume"`
	Capacity     string `json:"capacity"`
	AccessModes  string `json:"accessModes"`
	StorageClass string `json:"storageClass"`
	Age          string `json:"age"`
}

type StorageClassItem struct {
	Name                 string `json:"name"`
	Provisioner          string `json:"provisioner"`
	ReclaimPolicy        string `json:"reclaimPolicy"`
	VolumeBindingMode    string `json:"volumeBindingMode"`
	AllowVolumeExpansion bool   `json:"allowVolumeExpansion"`
	Age                  string `json:"age"`
}

type NodeItem struct {
	Name        string `json:"name"`
	Roles       string `json:"roles"`
	Version     string `json:"version"`
	InternalIP  string `json:"internalIP"`
	Status      string `json:"status"`
	OSImage     string `json:"osImage"`
	CPU         string `json:"cpu"`
	Memory      string `json:"memory"`
	PodCount    int    `json:"podCount"`
	LabelsCount int    `json:"labelsCount"`
	TaintsCount int    `json:"taintsCount"`
	Age         string `json:"age"`
}

type LimitRangeItem struct {
	Name          string `json:"name"`
	Namespace     string `json:"namespace"`
	LimitsCount   int    `json:"limitsCount"`
	DefaultCPU    string `json:"defaultCpu"`
	DefaultMemory string `json:"defaultMemory"`
	Age           string `json:"age"`
}

type ResourceQuotaItem struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	HardPods   string `json:"hardPods"`
	UsedPods   string `json:"usedPods"`
	HardCPU    string `json:"hardCpu"`
	UsedCPU    string `json:"usedCpu"`
	HardMemory string `json:"hardMemory"`
	UsedMemory string `json:"usedMemory"`
	HardPVCs   string `json:"hardPvcs"`
	UsedPVCs   string `json:"usedPvcs"`
	Age        string `json:"age"`
}

type NetworkPolicyItem struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	PodSelector  string `json:"podSelector"`
	PolicyTypes  string `json:"policyTypes"`
	IngressRules int    `json:"ingressRules"`
	EgressRules  int    `json:"egressRules"`
	Age          string `json:"age"`
}

type ResourceService struct {
	services          []ServiceItem
	configMaps        []ConfigMapItem
	secrets           []SecretItem
	ingresses         []IngressItem
	hpas              []HPAItem
	pvs               []PVItem
	pvcs              []PVCItem
	scs               []StorageClassItem
	nodes             []NodeItem
	nodeYAML          map[string]string
	limitRanges       []LimitRangeItem
	resourceQuotas    []ResourceQuotaItem
	networkPolicies   []NetworkPolicyItem
	limitRangeYAML    map[string]string
	resourceQuotaYAML map[string]string
	networkPolicyYAML map[string]string
}

func NewResourceService() *ResourceService {
	return &ResourceService{
		services: []ServiceItem{
			{Name: "web-api-svc", Namespace: "default", Type: "ClusterIP", ClusterIP: "10.43.12.7", Ports: "80/TCP", Pods: 3, Age: "12d"},
			{Name: "worker-svc", Namespace: "dev", Type: "ClusterIP", ClusterIP: "10.43.44.2", Ports: "8080/TCP", Pods: 2, Age: "1d"},
		},
		configMaps: []ConfigMapItem{
			{Name: "web-api-config", Namespace: "default", DataCount: 4, Age: "9d"},
			{Name: "worker-config", Namespace: "dev", DataCount: 2, Age: "1d"},
		},
		secrets: []SecretItem{
			{
				Name:      "web-api-secret",
				Namespace: "default",
				Type:      "Opaque",
				Age:       "9d",
				Data: map[string]string{
					"DB_PASSWORD": maskSecret("super-secret-password"),
					"API_KEY":     maskSecret("abcd1234efgh5678"),
				},
			},
			{
				Name:      "worker-secret",
				Namespace: "dev",
				Type:      "Opaque",
				Age:       "1d",
				Data: map[string]string{
					"REDIS_PASSWORD": maskSecret("redis-pass"),
				},
			},
		},
		ingresses: []IngressItem{
			{
				Name:      "web-api-ing",
				Namespace: "default",
				ClassName: "nginx",
				Hosts:     []string{"api.example.com"},
				Address:   "10.0.0.88",
				TLS:       true,
				Age:       "6d",
			},
			{
				Name:      "worker-ing",
				Namespace: "dev",
				ClassName: "nginx",
				Hosts:     []string{"worker.example.com"},
				Address:   "10.0.0.89",
				TLS:       false,
				Age:       "1d",
			},
		},
		hpas: []HPAItem{
			{
				Name:              "web-api-hpa",
				Namespace:         "default",
				TargetKind:        "Deployment",
				TargetName:        "web-api",
				MinReplicas:       2,
				MaxReplicas:       6,
				CurrentReplicas:   3,
				TargetCPUPercent:  70,
				CurrentCPUPercent: 52,
				Age:               "5d",
			},
			{
				Name:              "worker-hpa",
				Namespace:         "dev",
				TargetKind:        "Deployment",
				TargetName:        "task-worker",
				MinReplicas:       1,
				MaxReplicas:       4,
				CurrentReplicas:   2,
				TargetCPUPercent:  75,
				CurrentCPUPercent: 63,
				Age:               "20h",
			},
		},
		pvs: []PVItem{
			{
				Name:          "pv-web-data",
				Capacity:      "20Gi",
				AccessModes:   "RWO",
				ReclaimPolicy: "Delete",
				Status:        "Bound",
				ClaimRef:      "default/web-data",
				StorageClass:  "fast-ssd",
				Age:           "14d",
			},
			{
				Name:          "pv-worker-cache",
				Capacity:      "10Gi",
				AccessModes:   "RWO",
				ReclaimPolicy: "Retain",
				Status:        "Bound",
				ClaimRef:      "dev/worker-cache",
				StorageClass:  "standard",
				Age:           "5d",
			},
		},
		pvcs: []PVCItem{
			{
				Name:         "web-data",
				Namespace:    "default",
				Status:       "Bound",
				Volume:       "pv-web-data",
				Capacity:     "20Gi",
				AccessModes:  "RWO",
				StorageClass: "fast-ssd",
				Age:          "14d",
			},
			{
				Name:         "worker-cache",
				Namespace:    "dev",
				Status:       "Bound",
				Volume:       "pv-worker-cache",
				Capacity:     "10Gi",
				AccessModes:  "RWO",
				StorageClass: "standard",
				Age:          "5d",
			},
		},
		scs: []StorageClassItem{
			{
				Name:                 "fast-ssd",
				Provisioner:          "kubernetes.io/no-provisioner",
				ReclaimPolicy:        "Delete",
				VolumeBindingMode:    "WaitForFirstConsumer",
				AllowVolumeExpansion: true,
				Age:                  "20d",
			},
			{
				Name:                 "standard",
				Provisioner:          "kubernetes.io/no-provisioner",
				ReclaimPolicy:        "Retain",
				VolumeBindingMode:    "Immediate",
				AllowVolumeExpansion: false,
				Age:                  "30d",
			},
		},
		nodes: []NodeItem{
			{
				Name:        "ip-10-10-1-21.ec2.internal",
				Roles:       "control-plane,master",
				Version:     "v1.30.2",
				InternalIP:  "10.10.1.21",
				Status:      "Ready",
				OSImage:     "Ubuntu 22.04.4 LTS",
				CPU:         "4",
				Memory:      "15728640Ki",
				PodCount:    42,
				LabelsCount: 18,
				TaintsCount: 1,
				Age:         "21d",
			},
			{
				Name:        "ip-10-10-1-35.ec2.internal",
				Roles:       "worker",
				Version:     "v1.30.2",
				InternalIP:  "10.10.1.35",
				Status:      "Ready",
				OSImage:     "Ubuntu 22.04.4 LTS",
				CPU:         "8",
				Memory:      "31457280Ki",
				PodCount:    89,
				LabelsCount: 14,
				TaintsCount: 0,
				Age:         "20d",
			},
		},
		nodeYAML: map[string]string{
			"ip-10-10-1-21.ec2.internal": "apiVersion: v1\nkind: Node\nmetadata:\n  name: ip-10-10-1-21.ec2.internal\n  labels:\n    node-role.kubernetes.io/control-plane: \"\"\nstatus:\n  nodeInfo:\n    kubeletVersion: v1.30.2\n",
			"ip-10-10-1-35.ec2.internal": "apiVersion: v1\nkind: Node\nmetadata:\n  name: ip-10-10-1-35.ec2.internal\n  labels:\n    node-role.kubernetes.io/worker: \"\"\nstatus:\n  nodeInfo:\n    kubeletVersion: v1.30.2\n",
		},
		limitRanges: []LimitRangeItem{
			{
				Name:          "compute-defaults",
				Namespace:     "default",
				LimitsCount:   1,
				DefaultCPU:    "500m",
				DefaultMemory: "512Mi",
				Age:           "12d",
			},
			{
				Name:          "dev-container-limits",
				Namespace:     "dev",
				LimitsCount:   1,
				DefaultCPU:    "300m",
				DefaultMemory: "256Mi",
				Age:           "3d",
			},
		},
		resourceQuotas: []ResourceQuotaItem{
			{
				Name:       "compute-quota",
				Namespace:  "default",
				HardPods:   "20",
				UsedPods:   "7",
				HardCPU:    "8",
				UsedCPU:    "2",
				HardMemory: "16Gi",
				UsedMemory: "4Gi",
				HardPVCs:   "10",
				UsedPVCs:   "3",
				Age:        "10d",
			},
			{
				Name:       "dev-quota",
				Namespace:  "dev",
				HardPods:   "15",
				UsedPods:   "5",
				HardCPU:    "4",
				UsedCPU:    "1",
				HardMemory: "8Gi",
				UsedMemory: "2Gi",
				HardPVCs:   "8",
				UsedPVCs:   "2",
				Age:        "4d",
			},
		},
		networkPolicies: []NetworkPolicyItem{
			{
				Name:         "default-deny-all",
				Namespace:    "default",
				PodSelector:  "<all>",
				PolicyTypes:  "Ingress,Egress",
				IngressRules: 0,
				EgressRules:  0,
				Age:          "9d",
			},
			{
				Name:         "allow-web-to-api",
				Namespace:    "dev",
				PodSelector:  "app=api",
				PolicyTypes:  "Ingress",
				IngressRules: 1,
				EgressRules:  0,
				Age:          "2d",
			},
		},
		limitRangeYAML: map[string]string{
			"compute-defaults":     "apiVersion: v1\nkind: LimitRange\nmetadata:\n  name: compute-defaults\n  namespace: default\nspec:\n  limits:\n  - type: Container\n    default:\n      cpu: 500m\n      memory: 512Mi\n    defaultRequest:\n      cpu: 100m\n      memory: 128Mi\n",
			"dev-container-limits": "apiVersion: v1\nkind: LimitRange\nmetadata:\n  name: dev-container-limits\n  namespace: dev\nspec:\n  limits:\n  - type: Container\n    default:\n      cpu: 300m\n      memory: 256Mi\n    defaultRequest:\n      cpu: 100m\n      memory: 128Mi\n",
		},
		resourceQuotaYAML: map[string]string{
			"compute-quota": "apiVersion: v1\nkind: ResourceQuota\nmetadata:\n  name: compute-quota\n  namespace: default\nspec:\n  hard:\n    pods: \"20\"\n    requests.cpu: \"8\"\n    requests.memory: 16Gi\n    persistentvolumeclaims: \"10\"\n",
			"dev-quota":     "apiVersion: v1\nkind: ResourceQuota\nmetadata:\n  name: dev-quota\n  namespace: dev\nspec:\n  hard:\n    pods: \"15\"\n    requests.cpu: \"4\"\n    requests.memory: 8Gi\n    persistentvolumeclaims: \"8\"\n",
		},
		networkPolicyYAML: map[string]string{
			"default-deny-all": "apiVersion: networking.k8s.io/v1\nkind: NetworkPolicy\nmetadata:\n  name: default-deny-all\n  namespace: default\nspec:\n  podSelector: {}\n  policyTypes:\n  - Ingress\n  - Egress\n",
			"allow-web-to-api": "apiVersion: networking.k8s.io/v1\nkind: NetworkPolicy\nmetadata:\n  name: allow-web-to-api\n  namespace: dev\nspec:\n  podSelector:\n    matchLabels:\n      app: api\n  policyTypes:\n  - Ingress\n  ingress:\n  - from:\n    - podSelector:\n        matchLabels:\n          app: web\n",
		},
	}
}

func (s *ResourceService) ListServices() []ServiceItem {
	return append([]ServiceItem(nil), s.services...)
}

func (s *ResourceService) GetService(name string) (ServiceItem, bool) {
	for _, item := range s.services {
		if item.Name == name {
			return item, true
		}
	}
	return ServiceItem{}, false
}

func (s *ResourceService) ListConfigMaps() []ConfigMapItem {
	return append([]ConfigMapItem(nil), s.configMaps...)
}

func (s *ResourceService) GetConfigMap(name string) (ConfigMapItem, bool) {
	for _, item := range s.configMaps {
		if item.Name == name {
			return item, true
		}
	}
	return ConfigMapItem{}, false
}

func (s *ResourceService) ListSecrets() []SecretItem {
	return append([]SecretItem(nil), s.secrets...)
}

func (s *ResourceService) GetSecret(name string) (SecretItem, bool) {
	for _, item := range s.secrets {
		if item.Name == name {
			return item, true
		}
	}
	return SecretItem{}, false
}

func (s *ResourceService) ListIngresses() []IngressItem {
	return append([]IngressItem(nil), s.ingresses...)
}

func (s *ResourceService) GetIngress(name string) (IngressItem, bool) {
	for _, item := range s.ingresses {
		if item.Name == name {
			return item, true
		}
	}
	return IngressItem{}, false
}

func (s *ResourceService) ListHPAs() []HPAItem {
	return append([]HPAItem(nil), s.hpas...)
}

func (s *ResourceService) GetHPA(name string) (HPAItem, bool) {
	for _, item := range s.hpas {
		if item.Name == name {
			return item, true
		}
	}
	return HPAItem{}, false
}

func (s *ResourceService) ListIngressServices(name string) ([]ServiceItem, bool) {
	switch name {
	case "web-api-ing":
		item, ok := s.GetService("web-api-svc")
		if !ok {
			return nil, false
		}
		return []ServiceItem{item}, true
	case "worker-ing":
		item, ok := s.GetService("worker-svc")
		if !ok {
			return nil, false
		}
		return []ServiceItem{item}, true
	default:
		return nil, false
	}
}

func (s *ResourceService) GetHPATarget(name string) (HPATarget, bool) {
	hpa, ok := s.GetHPA(name)
	if !ok {
		return HPATarget{}, false
	}
	return HPATarget{
		Kind:            hpa.TargetKind,
		Name:            hpa.TargetName,
		Namespace:       hpa.Namespace,
		CurrentReplicas: hpa.CurrentReplicas,
		DesiredReplicas: hpa.CurrentReplicas,
	}, true
}

func (s *ResourceService) ListPVs() []PVItem {
	return append([]PVItem(nil), s.pvs...)
}

func (s *ResourceService) GetPV(name string) (PVItem, bool) {
	for _, item := range s.pvs {
		if item.Name == name {
			return item, true
		}
	}
	return PVItem{}, false
}

func (s *ResourceService) ListPVCs() []PVCItem {
	return append([]PVCItem(nil), s.pvcs...)
}

func (s *ResourceService) GetPVC(name string) (PVCItem, bool) {
	for _, item := range s.pvcs {
		if item.Name == name {
			return item, true
		}
	}
	return PVCItem{}, false
}

func (s *ResourceService) ListStorageClasses() []StorageClassItem {
	return append([]StorageClassItem(nil), s.scs...)
}

func (s *ResourceService) GetStorageClass(name string) (StorageClassItem, bool) {
	for _, item := range s.scs {
		if item.Name == name {
			return item, true
		}
	}
	return StorageClassItem{}, false
}

func (s *ResourceService) ListNodes() []NodeItem {
	return append([]NodeItem(nil), s.nodes...)
}

func (s *ResourceService) GetNode(name string) (NodeItem, bool) {
	for _, item := range s.nodes {
		if item.Name == name {
			return item, true
		}
	}
	return NodeItem{}, false
}

func (s *ResourceService) GetNodeYAML(name string) (string, bool) {
	yaml, ok := s.nodeYAML[name]
	return yaml, ok
}

func (s *ResourceService) ListLimitRanges() []LimitRangeItem {
	return append([]LimitRangeItem(nil), s.limitRanges...)
}

func (s *ResourceService) GetLimitRange(name string) (LimitRangeItem, bool) {
	for _, item := range s.limitRanges {
		if item.Name == name {
			return item, true
		}
	}
	return LimitRangeItem{}, false
}

func (s *ResourceService) GetLimitRangeYAML(name string) (string, bool) {
	yaml, ok := s.limitRangeYAML[name]
	return yaml, ok
}

func (s *ResourceService) LimitRangeNamespace(name string) (string, error) {
	item, ok := s.GetLimitRange(name)
	if !ok {
		return "", fmt.Errorf("limitrange not found: %s", name)
	}
	return item.Namespace, nil
}

func (s *ResourceService) UpdateLimitRangeYAML(name, yaml string) error {
	return updateResourceYAML(s.limitRangeYAML, "limitrange", name, yaml)
}

func (s *ResourceService) CreateLimitRange(namespace, rawYAML string) error {
	name, err := extractMetadataNameFromYAML(rawYAML)
	if err != nil {
		return err
	}
	if strings.TrimSpace(namespace) == "" {
		return fmt.Errorf("namespace is required")
	}
	if _, exists := s.limitRangeYAML[name]; exists {
		return fmt.Errorf("limitrange already exists: %s", name)
	}
	s.limitRangeYAML[name] = rawYAML
	s.limitRanges = append(s.limitRanges, LimitRangeItem{
		Name:          name,
		Namespace:     strings.TrimSpace(namespace),
		LimitsCount:   1,
		DefaultCPU:    "",
		DefaultMemory: "",
		Age:           "0m",
	})
	return nil
}

func (s *ResourceService) DeleteLimitRange(name string) error {
	if _, exists := s.limitRangeYAML[name]; !exists {
		return fmt.Errorf("limitrange not found: %s", name)
	}
	delete(s.limitRangeYAML, name)
	for i, item := range s.limitRanges {
		if item.Name == name {
			s.limitRanges = append(s.limitRanges[:i], s.limitRanges[i+1:]...)
			return nil
		}
	}
	return nil
}

func (s *ResourceService) ListResourceQuotas() []ResourceQuotaItem {
	return append([]ResourceQuotaItem(nil), s.resourceQuotas...)
}

func (s *ResourceService) GetResourceQuota(name string) (ResourceQuotaItem, bool) {
	for _, item := range s.resourceQuotas {
		if item.Name == name {
			return item, true
		}
	}
	return ResourceQuotaItem{}, false
}

func (s *ResourceService) GetResourceQuotaYAML(name string) (string, bool) {
	yaml, ok := s.resourceQuotaYAML[name]
	return yaml, ok
}

func (s *ResourceService) ResourceQuotaNamespace(name string) (string, error) {
	item, ok := s.GetResourceQuota(name)
	if !ok {
		return "", fmt.Errorf("resourcequota not found: %s", name)
	}
	return item.Namespace, nil
}

func (s *ResourceService) UpdateResourceQuotaYAML(name, yaml string) error {
	return updateResourceYAML(s.resourceQuotaYAML, "resourcequota", name, yaml)
}

func (s *ResourceService) CreateResourceQuota(namespace, rawYAML string) error {
	name, err := extractMetadataNameFromYAML(rawYAML)
	if err != nil {
		return err
	}
	if strings.TrimSpace(namespace) == "" {
		return fmt.Errorf("namespace is required")
	}
	if _, exists := s.resourceQuotaYAML[name]; exists {
		return fmt.Errorf("resourcequota already exists: %s", name)
	}
	s.resourceQuotaYAML[name] = rawYAML
	s.resourceQuotas = append(s.resourceQuotas, ResourceQuotaItem{
		Name:       name,
		Namespace:  strings.TrimSpace(namespace),
		HardPods:   "",
		UsedPods:   "",
		HardCPU:    "",
		UsedCPU:    "",
		HardMemory: "",
		UsedMemory: "",
		HardPVCs:   "",
		UsedPVCs:   "",
		Age:        "0m",
	})
	return nil
}

func (s *ResourceService) DeleteResourceQuota(name string) error {
	if _, exists := s.resourceQuotaYAML[name]; !exists {
		return fmt.Errorf("resourcequota not found: %s", name)
	}
	delete(s.resourceQuotaYAML, name)
	for i, item := range s.resourceQuotas {
		if item.Name == name {
			s.resourceQuotas = append(s.resourceQuotas[:i], s.resourceQuotas[i+1:]...)
			return nil
		}
	}
	return nil
}

func (s *ResourceService) ListNetworkPolicies() []NetworkPolicyItem {
	return append([]NetworkPolicyItem(nil), s.networkPolicies...)
}

func (s *ResourceService) GetNetworkPolicy(name string) (NetworkPolicyItem, bool) {
	for _, item := range s.networkPolicies {
		if item.Name == name {
			return item, true
		}
	}
	return NetworkPolicyItem{}, false
}

func (s *ResourceService) GetNetworkPolicyYAML(name string) (string, bool) {
	yaml, ok := s.networkPolicyYAML[name]
	return yaml, ok
}

func (s *ResourceService) NetworkPolicyNamespace(name string) (string, error) {
	item, ok := s.GetNetworkPolicy(name)
	if !ok {
		return "", fmt.Errorf("networkpolicy not found: %s", name)
	}
	return item.Namespace, nil
}

func (s *ResourceService) UpdateNetworkPolicyYAML(name, yaml string) error {
	return updateResourceYAML(s.networkPolicyYAML, "networkpolicy", name, yaml)
}

func (s *ResourceService) CreateNetworkPolicy(namespace, rawYAML string) error {
	name, err := extractMetadataNameFromYAML(rawYAML)
	if err != nil {
		return err
	}
	if strings.TrimSpace(namespace) == "" {
		return fmt.Errorf("namespace is required")
	}
	if _, exists := s.networkPolicyYAML[name]; exists {
		return fmt.Errorf("networkpolicy already exists: %s", name)
	}
	s.networkPolicyYAML[name] = rawYAML
	s.networkPolicies = append(s.networkPolicies, NetworkPolicyItem{
		Name:         name,
		Namespace:    strings.TrimSpace(namespace),
		PodSelector:  "<all>",
		PolicyTypes:  "Ingress",
		IngressRules: 0,
		EgressRules:  0,
		Age:          "0m",
	})
	return nil
}

func (s *ResourceService) DeleteNetworkPolicy(name string) error {
	if _, exists := s.networkPolicyYAML[name]; !exists {
		return fmt.Errorf("networkpolicy not found: %s", name)
	}
	delete(s.networkPolicyYAML, name)
	for i, item := range s.networkPolicies {
		if item.Name == name {
			s.networkPolicies = append(s.networkPolicies[:i], s.networkPolicies[i+1:]...)
			return nil
		}
	}
	return nil
}

func updateResourceYAML(target map[string]string, kind, name, yaml string) error {
	if strings.TrimSpace(yaml) == "" {
		return fmt.Errorf("yaml content is empty")
	}
	if _, exists := target[name]; !exists {
		return fmt.Errorf("%s not found: %s", kind, name)
	}
	target[name] = yaml
	return nil
}

func extractMetadataNameFromYAML(raw string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", fmt.Errorf("yaml content is empty")
	}
	var obj map[string]any
	if err := yaml.Unmarshal([]byte(raw), &obj); err != nil {
		return "", fmt.Errorf("invalid yaml: %w", err)
	}
	metadata, ok := obj["metadata"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("yaml metadata.name is required")
	}
	name, _ := metadata["name"].(string)
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("yaml metadata.name is required")
	}
	return name, nil
}

func maskSecret(value string) string {
	if len(value) <= 4 {
		return "****"
	}
	return value[:2] + "******" + value[len(value)-2:]
}
