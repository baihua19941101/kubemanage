package handlers

import (
	"net/http"

	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type NamespaceHandler struct {
	namespaceSvc         *service.NamespaceService
	clusterConnectionSvc *service.ClusterConnectionService
	adapterMode          string
}

type CreateNamespaceRequest struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

func NewNamespaceHandler(namespaceSvc *service.NamespaceService, clusterConnectionSvc *service.ClusterConnectionService, adapterMode string) *NamespaceHandler {
	return &NamespaceHandler{
		namespaceSvc:         namespaceSvc,
		clusterConnectionSvc: clusterConnectionSvc,
		adapterMode:          adapterMode,
	}
}

func (h *NamespaceHandler) ListNamespaces(c *gin.Context) {
	if h.adapterMode != "mock" && h.clusterConnectionSvc != nil {
		items, err := h.clusterConnectionSvc.ListLiveNamespaces(c.Request.Context())
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"items": items})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.namespaceSvc.List()})
}

func (h *NamespaceHandler) CreateNamespace(c *gin.Context) {
	var req CreateNamespaceRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if h.adapterMode != "mock" {
		if h.clusterConnectionSvc == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "cluster connection service unavailable"})
			return
		}
		ns, err := h.clusterConnectionSvc.CreateLiveNamespace(c.Request.Context(), req.Name, req.Labels)
		if err != nil {
			if err.Error() == "namespace name is required" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "namespace already exists: "+req.Name {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, ns)
		return
	}

	ns, err := h.namespaceSvc.Create(req.Name, req.Labels)
	if err != nil {
		if err.Error() == "namespace name is required" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "namespace already exists: "+req.Name {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ns)
}

func (h *NamespaceHandler) DeleteNamespace(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "namespace name is required"})
		return
	}
	if h.adapterMode != "mock" {
		if h.clusterConnectionSvc == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "cluster connection service unavailable"})
			return
		}
		if err := h.clusterConnectionSvc.DeleteLiveNamespace(c.Request.Context(), name); err != nil {
			if err.Error() == "namespace name is required" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "namespace not found: "+name {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
		return
	}

	if err := h.namespaceSvc.Delete(name); err != nil {
		if err.Error() == "namespace not found: "+name {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *NamespaceHandler) GetNamespaceYAML(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" {
		if h.clusterConnectionSvc == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "cluster connection service unavailable"})
			return
		}
		content, err := h.clusterConnectionSvc.GetLiveNamespaceYAML(c.Request.Context(), name)
		if err != nil {
			if err.Error() == "namespace not found: "+name {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/yaml; charset=utf-8")
		c.String(http.StatusOK, content)
		return
	}
	content, err := h.namespaceSvc.YAML(name)
	if err != nil {
		if err.Error() == "namespace not found: "+name {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.String(http.StatusOK, content)
}

func (h *NamespaceHandler) DownloadNamespaceYAML(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" {
		if h.clusterConnectionSvc == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "cluster connection service unavailable"})
			return
		}
		content, err := h.clusterConnectionSvc.GetLiveNamespaceYAML(c.Request.Context(), name)
		if err != nil {
			if err.Error() == "namespace not found: "+name {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/yaml; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=\"namespace-"+name+".yaml\"")
		c.String(http.StatusOK, content)
		return
	}
	content, err := h.namespaceSvc.YAML(name)
	if err != nil {
		if err.Error() == "namespace not found: "+name {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=\"namespace-"+name+".yaml\"")
	c.String(http.StatusOK, content)
}
