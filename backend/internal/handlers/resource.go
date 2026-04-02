package handlers

import (
	"net/http"
	"strings"

	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ResourceHandler struct {
	resourceSvc     *service.ResourceService
	liveResourceSvc *service.LiveResourceReader
	adapterMode     string
}

type updateResourceYAMLRequest struct {
	YAML string `json:"yaml"`
}

type createResourceYAMLRequest struct {
	Namespace string `json:"namespace"`
	YAML      string `json:"yaml"`
}

func NewResourceHandler(resourceSvc *service.ResourceService, liveResourceSvc *service.LiveResourceReader, adapterMode string) *ResourceHandler {
	return &ResourceHandler{
		resourceSvc:     resourceSvc,
		liveResourceSvc: liveResourceSvc,
		adapterMode:     adapterMode,
	}
}

func (h *ResourceHandler) ListServices(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		items, err := h.liveResourceSvc.ListServices(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListServices()})
}

func (h *ResourceHandler) GetService(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		item, err := h.liveResourceSvc.GetService(c.Request.Context(), c.Param("name"))
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
		return
	}
	item, ok := h.resourceSvc.GetService(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListConfigMaps(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		items, err := h.liveResourceSvc.ListConfigMaps(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListConfigMaps()})
}

func (h *ResourceHandler) GetConfigMap(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		item, err := h.liveResourceSvc.GetConfigMap(c.Request.Context(), c.Param("name"))
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
		return
	}
	item, ok := h.resourceSvc.GetConfigMap(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "configmap not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListSecrets(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		items, err := h.liveResourceSvc.ListSecrets(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListSecrets()})
}

func (h *ResourceHandler) GetSecret(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		item, err := h.liveResourceSvc.GetSecret(c.Request.Context(), c.Param("name"))
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
		return
	}
	item, ok := h.resourceSvc.GetSecret(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "secret not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListIngresses(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		items, err := h.liveResourceSvc.ListIngresses(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListIngresses()})
}

func (h *ResourceHandler) GetIngress(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		item, err := h.liveResourceSvc.GetIngress(c.Request.Context(), c.Param("name"))
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
		return
	}
	item, ok := h.resourceSvc.GetIngress(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "ingress not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListHPAs(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		items, err := h.liveResourceSvc.ListHPAs(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListHPAs()})
}

func (h *ResourceHandler) GetHPA(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		item, err := h.liveResourceSvc.GetHPA(c.Request.Context(), c.Param("name"))
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
		return
	}
	item, ok := h.resourceSvc.GetHPA(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "hpa not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListIngressServices(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		items, err := h.liveResourceSvc.ListIngressServices(c.Request.Context(), c.Param("name"))
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	items, ok := h.resourceSvc.ListIngressServices(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "ingress not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ResourceHandler) GetHPATarget(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		item, err := h.liveResourceSvc.GetHPATarget(c.Request.Context(), c.Param("name"))
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
		return
	}
	item, ok := h.resourceSvc.GetHPATarget(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "hpa not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListPVs(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		items, err := h.liveResourceSvc.ListPVs(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListPVs()})
}

func (h *ResourceHandler) GetPV(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		item, err := h.liveResourceSvc.GetPV(c.Request.Context(), c.Param("name"))
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
		return
	}
	item, ok := h.resourceSvc.GetPV(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "pv not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListPVCs(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		items, err := h.liveResourceSvc.ListPVCs(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListPVCs()})
}

func (h *ResourceHandler) GetPVC(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		item, err := h.liveResourceSvc.GetPVC(c.Request.Context(), c.Param("name"))
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
		return
	}
	item, ok := h.resourceSvc.GetPVC(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "pvc not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListStorageClasses(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		items, err := h.liveResourceSvc.ListStorageClasses(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListStorageClasses()})
}

func (h *ResourceHandler) GetStorageClass(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		item, err := h.liveResourceSvc.GetStorageClass(c.Request.Context(), c.Param("name"))
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
		return
	}
	item, ok := h.resourceSvc.GetStorageClass(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "storageclass not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListNodes(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		items, err := h.liveResourceSvc.ListNodes(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListNodes()})
}

func (h *ResourceHandler) GetNode(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		item, err := h.liveResourceSvc.GetNode(c.Request.Context(), name)
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
		return
	}
	item, ok := h.resourceSvc.GetNode(name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) GetNodeYAML(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		raw, err := h.liveResourceSvc.GetNodeYAML(c.Request.Context(), name)
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/yaml; charset=utf-8")
		c.String(http.StatusOK, raw)
		return
	}
	raw, ok := h.resourceSvc.GetNodeYAML(name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.String(http.StatusOK, raw)
}

func (h *ResourceHandler) DownloadNodeYAML(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		raw, err := h.liveResourceSvc.GetNodeYAML(c.Request.Context(), name)
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/yaml; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=\"node-"+name+".yaml\"")
		c.String(http.StatusOK, raw)
		return
	}
	raw, ok := h.resourceSvc.GetNodeYAML(name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=\"node-"+name+".yaml\"")
	c.String(http.StatusOK, raw)
}

func (h *ResourceHandler) ListLimitRanges(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		items, err := h.liveResourceSvc.ListLimitRanges(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListLimitRanges()})
}

func (h *ResourceHandler) GetLimitRange(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		item, err := h.liveResourceSvc.GetLimitRange(c.Request.Context(), name)
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
		return
	}
	item, ok := h.resourceSvc.GetLimitRange(name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "limitrange not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) GetLimitRangeYAML(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		raw, err := h.liveResourceSvc.GetLimitRangeYAML(c.Request.Context(), name)
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/yaml; charset=utf-8")
		c.String(http.StatusOK, raw)
		return
	}
	raw, ok := h.resourceSvc.GetLimitRangeYAML(name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "limitrange not found"})
		return
	}
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.String(http.StatusOK, raw)
}

func (h *ResourceHandler) DownloadLimitRangeYAML(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		raw, err := h.liveResourceSvc.GetLimitRangeYAML(c.Request.Context(), name)
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/yaml; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=\"limitrange-"+name+".yaml\"")
		c.String(http.StatusOK, raw)
		return
	}
	raw, ok := h.resourceSvc.GetLimitRangeYAML(name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "limitrange not found"})
		return
	}
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=\"limitrange-"+name+".yaml\"")
	c.String(http.StatusOK, raw)
}

func (h *ResourceHandler) ListResourceQuotas(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		items, err := h.liveResourceSvc.ListResourceQuotas(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListResourceQuotas()})
}

func (h *ResourceHandler) GetResourceQuota(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		item, err := h.liveResourceSvc.GetResourceQuota(c.Request.Context(), name)
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
		return
	}
	item, ok := h.resourceSvc.GetResourceQuota(name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "resourcequota not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) GetResourceQuotaYAML(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		raw, err := h.liveResourceSvc.GetResourceQuotaYAML(c.Request.Context(), name)
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/yaml; charset=utf-8")
		c.String(http.StatusOK, raw)
		return
	}
	raw, ok := h.resourceSvc.GetResourceQuotaYAML(name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "resourcequota not found"})
		return
	}
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.String(http.StatusOK, raw)
}

func (h *ResourceHandler) DownloadResourceQuotaYAML(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		raw, err := h.liveResourceSvc.GetResourceQuotaYAML(c.Request.Context(), name)
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/yaml; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=\"resourcequota-"+name+".yaml\"")
		c.String(http.StatusOK, raw)
		return
	}
	raw, ok := h.resourceSvc.GetResourceQuotaYAML(name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "resourcequota not found"})
		return
	}
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=\"resourcequota-"+name+".yaml\"")
	c.String(http.StatusOK, raw)
}

func (h *ResourceHandler) ListNetworkPolicies(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		items, err := h.liveResourceSvc.ListNetworkPolicies(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListNetworkPolicies()})
}

func (h *ResourceHandler) GetNetworkPolicy(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		item, err := h.liveResourceSvc.GetNetworkPolicy(c.Request.Context(), name)
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
		return
	}
	item, ok := h.resourceSvc.GetNetworkPolicy(name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "networkpolicy not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) GetNetworkPolicyYAML(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		raw, err := h.liveResourceSvc.GetNetworkPolicyYAML(c.Request.Context(), name)
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/yaml; charset=utf-8")
		c.String(http.StatusOK, raw)
		return
	}
	raw, ok := h.resourceSvc.GetNetworkPolicyYAML(name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "networkpolicy not found"})
		return
	}
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.String(http.StatusOK, raw)
}

func (h *ResourceHandler) DownloadNetworkPolicyYAML(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		raw, err := h.liveResourceSvc.GetNetworkPolicyYAML(c.Request.Context(), name)
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/yaml; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=\"networkpolicy-"+name+".yaml\"")
		c.String(http.StatusOK, raw)
		return
	}
	raw, ok := h.resourceSvc.GetNetworkPolicyYAML(name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "networkpolicy not found"})
		return
	}
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=\"networkpolicy-"+name+".yaml\"")
	c.String(http.StatusOK, raw)
}

func (h *ResourceHandler) UpdateLimitRangeYAML(c *gin.Context) {
	name := c.Param("name")
	var req updateResourceYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		if err := h.liveResourceSvc.UpdateLimitRangeYAML(c.Request.Context(), name, req.YAML); err != nil {
			handleResourceYAMLUpdateError(c, err, "limitrange")
			return
		}
		c.Status(http.StatusNoContent)
		return
	}
	if err := h.resourceSvc.UpdateLimitRangeYAML(name, req.YAML); err != nil {
		handleResourceYAMLUpdateError(c, err, "limitrange")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ResourceHandler) CreateLimitRange(c *gin.Context) {
	var req createResourceYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		if err := h.liveResourceSvc.CreateLimitRangeYAML(c.Request.Context(), req.Namespace, req.YAML); err != nil {
			handleResourceCreateError(c, err, "limitrange")
			return
		}
		c.Status(http.StatusCreated)
		return
	}
	if err := h.resourceSvc.CreateLimitRange(req.Namespace, req.YAML); err != nil {
		handleResourceCreateError(c, err, "limitrange")
		return
	}
	c.Status(http.StatusCreated)
}

func (h *ResourceHandler) UpdateResourceQuotaYAML(c *gin.Context) {
	name := c.Param("name")
	var req updateResourceYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		if err := h.liveResourceSvc.UpdateResourceQuotaYAML(c.Request.Context(), name, req.YAML); err != nil {
			handleResourceYAMLUpdateError(c, err, "resourcequota")
			return
		}
		c.Status(http.StatusNoContent)
		return
	}
	if err := h.resourceSvc.UpdateResourceQuotaYAML(name, req.YAML); err != nil {
		handleResourceYAMLUpdateError(c, err, "resourcequota")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ResourceHandler) CreateResourceQuota(c *gin.Context) {
	var req createResourceYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		if err := h.liveResourceSvc.CreateResourceQuotaYAML(c.Request.Context(), req.Namespace, req.YAML); err != nil {
			handleResourceCreateError(c, err, "resourcequota")
			return
		}
		c.Status(http.StatusCreated)
		return
	}
	if err := h.resourceSvc.CreateResourceQuota(req.Namespace, req.YAML); err != nil {
		handleResourceCreateError(c, err, "resourcequota")
		return
	}
	c.Status(http.StatusCreated)
}

func (h *ResourceHandler) UpdateNetworkPolicyYAML(c *gin.Context) {
	name := c.Param("name")
	var req updateResourceYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		if err := h.liveResourceSvc.UpdateNetworkPolicyYAML(c.Request.Context(), name, req.YAML); err != nil {
			handleResourceYAMLUpdateError(c, err, "networkpolicy")
			return
		}
		c.Status(http.StatusNoContent)
		return
	}
	if err := h.resourceSvc.UpdateNetworkPolicyYAML(name, req.YAML); err != nil {
		handleResourceYAMLUpdateError(c, err, "networkpolicy")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ResourceHandler) CreateNetworkPolicy(c *gin.Context) {
	var req createResourceYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		if err := h.liveResourceSvc.CreateNetworkPolicyYAML(c.Request.Context(), req.Namespace, req.YAML); err != nil {
			handleResourceCreateError(c, err, "networkpolicy")
			return
		}
		c.Status(http.StatusCreated)
		return
	}
	if err := h.resourceSvc.CreateNetworkPolicy(req.Namespace, req.YAML); err != nil {
		handleResourceCreateError(c, err, "networkpolicy")
		return
	}
	c.Status(http.StatusCreated)
}

func (h *ResourceHandler) DeleteLimitRange(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		if err := h.liveResourceSvc.DeleteLimitRange(c.Request.Context(), name); err != nil {
			handleResourceDeleteError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
		return
	}
	if err := h.resourceSvc.DeleteLimitRange(name); err != nil {
		handleResourceDeleteError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ResourceHandler) DeleteResourceQuota(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		if err := h.liveResourceSvc.DeleteResourceQuota(c.Request.Context(), name); err != nil {
			handleResourceDeleteError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
		return
	}
	if err := h.resourceSvc.DeleteResourceQuota(name); err != nil {
		handleResourceDeleteError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ResourceHandler) DeleteNetworkPolicy(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveResourceSvc != nil {
		if err := h.liveResourceSvc.DeleteNetworkPolicy(c.Request.Context(), name); err != nil {
			handleResourceDeleteError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
		return
	}
	if err := h.resourceSvc.DeleteNetworkPolicy(name); err != nil {
		handleResourceDeleteError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func handleResourceYAMLUpdateError(c *gin.Context, err error, kind string) {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "yaml content is empty"), strings.Contains(msg, "invalid yaml"):
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
	case strings.Contains(msg, "not found:"):
		c.JSON(http.StatusNotFound, gin.H{"error": msg})
	case strings.Contains(msg, "failed:"):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": msg})
	default:
		if strings.Contains(msg, kind+" not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": msg})
	}
}

func handleResourceDeleteError(c *gin.Context, err error) {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "not found:"):
		c.JSON(http.StatusNotFound, gin.H{"error": msg})
	case strings.Contains(msg, "failed:"):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": msg})
	default:
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": msg})
	}
}

func handleResourceCreateError(c *gin.Context, err error, kind string) {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "yaml content is empty"), strings.Contains(msg, "invalid yaml"), strings.Contains(msg, "namespace is required"), strings.Contains(msg, "metadata.name is required"):
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
	case strings.Contains(msg, "already exists"):
		c.JSON(http.StatusConflict, gin.H{"error": msg})
	case strings.Contains(msg, "failed:"):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": msg})
	default:
		if strings.Contains(msg, kind+" already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": msg})
	}
}
