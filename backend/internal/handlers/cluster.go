package handlers

import (
	"net/http"
	"strings"

	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ClusterHandler struct {
	clusterSvc           *service.ClusterService
	clusterConnectionSvc *service.ClusterConnectionService
	adapterMode          string
}

type SwitchClusterRequest struct {
	Name string `json:"name"`
}

func NewClusterHandler(clusterSvc *service.ClusterService, clusterConnectionSvc *service.ClusterConnectionService, adapterMode string) *ClusterHandler {
	return &ClusterHandler{
		clusterSvc:           clusterSvc,
		clusterConnectionSvc: clusterConnectionSvc,
		adapterMode:          adapterMode,
	}
}

func (h *ClusterHandler) ListClusters(c *gin.Context) {
	if live, ok, err := h.tryLiveCluster(c); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	} else if ok {
		c.JSON(http.StatusOK, gin.H{
			"items": []service.ClusterSummary{
				{
					State:             live.State,
					Name:              live.Name,
					Provider:          live.Provider,
					Distro:            live.Distro,
					KubernetesVersion: live.KubernetesVersion,
					Architecture:      live.Architecture,
					CPU:               live.CPU,
					Memory:            live.Memory,
					Pods:              live.Pods,
				},
			},
			"current": live.Name,
		})
		return
	}

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
	if live, ok, err := h.tryLiveCluster(c); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	} else if ok {
		c.JSON(http.StatusOK, service.ClusterSummary{
			State:             live.State,
			Name:              live.Name,
			Provider:          live.Provider,
			Distro:            live.Distro,
			KubernetesVersion: live.KubernetesVersion,
			Architecture:      live.Architecture,
			CPU:               live.CPU,
			Memory:            live.Memory,
			Pods:              live.Pods,
		})
		return
	}

	current, err := h.clusterSvc.GetCurrent(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, current)
}

func (h *ClusterHandler) tryLiveCluster(c *gin.Context) (service.LiveClusterSummary, bool, error) {
	if h.adapterMode == "mock" || h.clusterConnectionSvc == nil {
		return service.LiveClusterSummary{}, false, nil
	}
	live, err := h.clusterConnectionSvc.GetLiveCluster(c.Request.Context())
	if err == nil {
		return live, true, nil
	}
	return service.LiveClusterSummary{}, false, err
}

func (h *ClusterHandler) SwitchCluster(c *gin.Context) {
	var req SwitchClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	targetName := strings.TrimSpace(req.Name)

	if h.adapterMode != "mock" {
		if h.clusterConnectionSvc == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "cluster connection service unavailable"})
			return
		}
		if err := h.clusterConnectionSvc.ActivateByName(c.Request.Context(), targetName); err != nil {
			if err.Error() == "cluster name is required" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "cluster connection not found: "+targetName {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		current, err := h.clusterConnectionSvc.GetLiveCluster(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, current)
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
