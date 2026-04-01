package handlers

import (
	"net/http"
	"strconv"

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
	logs, err := h.workloadSvc.GetPodLogs(name, service.PodLogQuery{
		Keyword:       c.Query("keyword"),
		CaseSensitive: parseBool(c.Query("caseSensitive")),
		MatchOnly:     parseBool(c.Query("matchOnly")),
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, logs)
}

func (h *WorkloadHandler) GetTerminalCapabilities(c *gin.Context) {
	name := c.Param("name")
	caps, err := h.workloadSvc.GetTerminalCapabilities(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, caps)
}

func (h *WorkloadHandler) CreateTerminalSession(c *gin.Context) {
	name := c.Param("name")
	if err := h.workloadSvc.CreateTerminalSession(name); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "terminal gateway not enabled",
		"enabled": false,
	})
}

func (h *WorkloadHandler) ListStatefulSets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": h.workloadSvc.ListStatefulSets()})
}

func (h *WorkloadHandler) GetStatefulSet(c *gin.Context) {
	name := c.Param("name")
	item, err := h.workloadSvc.GetStatefulSet(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *WorkloadHandler) GetStatefulSetYAML(c *gin.Context) {
	name := c.Param("name")
	yaml, err := h.workloadSvc.GetStatefulSetYAML(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.String(http.StatusOK, yaml)
}

func (h *WorkloadHandler) UpdateStatefulSetYAML(c *gin.Context) {
	name := c.Param("name")
	var req UpdateYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.workloadSvc.UpdateStatefulSetYAML(name, req.YAML); err != nil {
		if err.Error() == "yaml content is empty" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *WorkloadHandler) ListDaemonSets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": h.workloadSvc.ListDaemonSets()})
}

func (h *WorkloadHandler) GetDaemonSet(c *gin.Context) {
	name := c.Param("name")
	item, err := h.workloadSvc.GetDaemonSet(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *WorkloadHandler) GetDaemonSetYAML(c *gin.Context) {
	name := c.Param("name")
	yaml, err := h.workloadSvc.GetDaemonSetYAML(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.String(http.StatusOK, yaml)
}

func (h *WorkloadHandler) UpdateDaemonSetYAML(c *gin.Context) {
	name := c.Param("name")
	var req UpdateYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.workloadSvc.UpdateDaemonSetYAML(name, req.YAML); err != nil {
		if err.Error() == "yaml content is empty" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *WorkloadHandler) ListJobs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": h.workloadSvc.ListJobs()})
}

func (h *WorkloadHandler) GetJob(c *gin.Context) {
	name := c.Param("name")
	item, err := h.workloadSvc.GetJob(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *WorkloadHandler) GetJobYAML(c *gin.Context) {
	name := c.Param("name")
	yaml, err := h.workloadSvc.GetJobYAML(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.String(http.StatusOK, yaml)
}

func (h *WorkloadHandler) UpdateJobYAML(c *gin.Context) {
	name := c.Param("name")
	var req UpdateYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.workloadSvc.UpdateJobYAML(name, req.YAML); err != nil {
		if err.Error() == "yaml content is empty" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *WorkloadHandler) ListCronJobs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": h.workloadSvc.ListCronJobs()})
}

func (h *WorkloadHandler) GetCronJob(c *gin.Context) {
	name := c.Param("name")
	item, err := h.workloadSvc.GetCronJob(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func parseBool(value string) bool {
	ok, err := strconv.ParseBool(value)
	return err == nil && ok
}

func (h *WorkloadHandler) GetCronJobYAML(c *gin.Context) {
	name := c.Param("name")
	yaml, err := h.workloadSvc.GetCronJobYAML(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.String(http.StatusOK, yaml)
}

func (h *WorkloadHandler) UpdateCronJobYAML(c *gin.Context) {
	name := c.Param("name")
	var req UpdateYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.workloadSvc.UpdateCronJobYAML(name, req.YAML); err != nil {
		if err.Error() == "yaml content is empty" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
