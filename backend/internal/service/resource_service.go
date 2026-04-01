package service

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
	Name               string `json:"name"`
	Namespace          string `json:"namespace"`
	TargetKind         string `json:"targetKind"`
	TargetName         string `json:"targetName"`
	MinReplicas        int    `json:"minReplicas"`
	MaxReplicas        int    `json:"maxReplicas"`
	CurrentReplicas    int    `json:"currentReplicas"`
	TargetCPUPercent   int    `json:"targetCPUPercent"`
	CurrentCPUPercent  int    `json:"currentCPUPercent"`
	Age                string `json:"age"`
}

type HPATarget struct {
	Kind            string `json:"kind"`
	Name            string `json:"name"`
	Namespace       string `json:"namespace"`
	CurrentReplicas int    `json:"currentReplicas"`
	DesiredReplicas int    `json:"desiredReplicas"`
}

type ResourceService struct {
	services   []ServiceItem
	configMaps []ConfigMapItem
	secrets    []SecretItem
	ingresses  []IngressItem
	hpas       []HPAItem
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

func maskSecret(value string) string {
	if len(value) <= 4 {
		return "****"
	}
	return value[:2] + "******" + value[len(value)-2:]
}
