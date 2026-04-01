package handlers

import (
	"net/http"

	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ResourceHandler struct {
	resourceSvc *service.ResourceService
}

func NewResourceHandler(resourceSvc *service.ResourceService) *ResourceHandler {
	return &ResourceHandler{resourceSvc: resourceSvc}
}

func (h *ResourceHandler) ListServices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListServices()})
}

func (h *ResourceHandler) GetService(c *gin.Context) {
	item, ok := h.resourceSvc.GetService(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListConfigMaps(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListConfigMaps()})
}

func (h *ResourceHandler) GetConfigMap(c *gin.Context) {
	item, ok := h.resourceSvc.GetConfigMap(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "configmap not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListSecrets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListSecrets()})
}

func (h *ResourceHandler) GetSecret(c *gin.Context) {
	item, ok := h.resourceSvc.GetSecret(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "secret not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListIngresses(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListIngresses()})
}

func (h *ResourceHandler) GetIngress(c *gin.Context) {
	item, ok := h.resourceSvc.GetIngress(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "ingress not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListHPAs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListHPAs()})
}

func (h *ResourceHandler) GetHPA(c *gin.Context) {
	item, ok := h.resourceSvc.GetHPA(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "hpa not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListIngressServices(c *gin.Context) {
	items, ok := h.resourceSvc.ListIngressServices(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "ingress not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ResourceHandler) GetHPATarget(c *gin.Context) {
	item, ok := h.resourceSvc.GetHPATarget(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "hpa not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListPVs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListPVs()})
}

func (h *ResourceHandler) GetPV(c *gin.Context) {
	item, ok := h.resourceSvc.GetPV(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "pv not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListPVCs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListPVCs()})
}

func (h *ResourceHandler) GetPVC(c *gin.Context) {
	item, ok := h.resourceSvc.GetPVC(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "pvc not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler) ListStorageClasses(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": h.resourceSvc.ListStorageClasses()})
}

func (h *ResourceHandler) GetStorageClass(c *gin.Context) {
	item, ok := h.resourceSvc.GetStorageClass(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "storageclass not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}
