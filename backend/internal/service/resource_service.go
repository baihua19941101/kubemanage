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

type ResourceService struct {
	services   []ServiceItem
	configMaps []ConfigMapItem
	secrets    []SecretItem
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

func maskSecret(value string) string {
	if len(value) <= 4 {
		return "****"
	}
	return value[:2] + "******" + value[len(value)-2:]
}
