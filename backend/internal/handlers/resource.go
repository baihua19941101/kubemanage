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
