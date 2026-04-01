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
