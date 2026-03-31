package handlers

import (
	"net/http"

	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type WorkloadHandler struct {
	workloadSvc *service.WorkloadService
}

type UpdateYAMLRequest struct {
	YAML string `json:"yaml"`
}

func NewWorkloadHandler(workloadSvc *service.WorkloadService) *WorkloadHandler {
	return &WorkloadHandler{workloadSvc: workloadSvc}
}

func (h *WorkloadHandler) ListDeployments(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": h.workloadSvc.ListDeployments()})
}

func (h *WorkloadHandler) GetDeployment(c *gin.Context) {
	name := c.Param("name")
	item, err := h.workloadSvc.GetDeployment(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *WorkloadHandler) GetDeploymentYAML(c *gin.Context) {
	name := c.Param("name")
	yaml, err := h.workloadSvc.GetDeploymentYAML(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.String(http.StatusOK, yaml)
}

func (h *WorkloadHandler) UpdateDeploymentYAML(c *gin.Context) {
	name := c.Param("name")
	var req UpdateYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.workloadSvc.UpdateDeploymentYAML(name, req.YAML); err != nil {
		if err.Error() == "yaml content is empty" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *WorkloadHandler) ListPods(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": h.workloadSvc.ListPods()})
}

func (h *WorkloadHandler) GetPod(c *gin.Context) {
	name := c.Param("name")
	item, err := h.workloadSvc.GetPod(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *WorkloadHandler) GetPodYAML(c *gin.Context) {
	name := c.Param("name")
	yaml, err := h.workloadSvc.GetPodYAML(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.String(http.StatusOK, yaml)
}

func (h *WorkloadHandler) UpdatePodYAML(c *gin.Context) {
	name := c.Param("name")
	var req UpdateYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.workloadSvc.UpdatePodYAML(name, req.YAML); err != nil {
		if err.Error() == "yaml content is empty" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *WorkloadHandler) GetPodLogs(c *gin.Context) {
	name := c.Param("name")
	logs, err := h.workloadSvc.GetPodLogs(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, logs)
}
