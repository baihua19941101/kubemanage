package handlers

import (
	"net/http"

	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ClusterHandler struct {
	clusterSvc *service.ClusterService
}

type SwitchClusterRequest struct {
	Name string `json:"name"`
}

func NewClusterHandler(clusterSvc *service.ClusterService) *ClusterHandler {
	return &ClusterHandler{
		clusterSvc: clusterSvc,
	}
}

func (h *ClusterHandler) ListClusters(c *gin.Context) {
	current, err := h.clusterSvc.GetCurrent(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":   h.clusterSvc.List(),
		"current": current.Name,
	})
}

func (h *ClusterHandler) GetCurrentCluster(c *gin.Context) {
	current, err := h.clusterSvc.GetCurrent(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, current)
}

func (h *ClusterHandler) SwitchCluster(c *gin.Context) {
	var req SwitchClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.clusterSvc.Switch(c.Request.Context(), req.Name); err != nil {
		if err.Error() == "cluster not found: "+req.Name {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	current, err := h.clusterSvc.GetCurrent(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, current)
}
