package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type WorkloadHandler struct {
	workloadSvc      *service.WorkloadService
	liveWorkloadSvc  *service.LiveWorkloadReader
	terminalSessions *service.TerminalSessionStore
	adapterMode      string
}

type UpdateYAMLRequest struct {
	YAML string `json:"yaml"`
}

type createTerminalSessionRequest struct {
	Container string `json:"container"`
}

func NewWorkloadHandler(workloadSvc *service.WorkloadService, liveWorkloadSvc *service.LiveWorkloadReader, terminalSessions *service.TerminalSessionStore, adapterMode string) *WorkloadHandler {
	return &WorkloadHandler{
		workloadSvc:      workloadSvc,
		liveWorkloadSvc:  liveWorkloadSvc,
		terminalSessions: terminalSessions,
		adapterMode:      adapterMode,
	}
}

func (h *WorkloadHandler) ListDeployments(c *gin.Context) {
	if h.adapterMode != "mock" && h.liveWorkloadSvc != nil {
		items, err := h.liveWorkloadSvc.ListDeployments(c.Request.Context())
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"items": items})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.workloadSvc.ListDeployments()})
}

func (h *WorkloadHandler) GetDeployment(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveWorkloadSvc != nil {
		item, err := h.liveWorkloadSvc.GetDeployment(c.Request.Context(), name)
		if err == nil {
			c.JSON(http.StatusOK, item)
			return
		}
		if strings.Contains(err.Error(), "not found:") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	item, err := h.workloadSvc.GetDeployment(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *WorkloadHandler) GetDeploymentYAML(c *gin.Context) {
	if h.adapterMode != "mock" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "deployment yaml path not enabled in real-only mode"})
		return
	}
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
	if h.adapterMode != "mock" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "deployment yaml write path not enabled in real-only mode"})
		return
	}
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
	if h.adapterMode != "mock" && h.liveWorkloadSvc != nil {
		items, err := h.liveWorkloadSvc.ListPods(c.Request.Context())
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"items": items})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.workloadSvc.ListPods()})
}

func (h *WorkloadHandler) GetPod(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveWorkloadSvc != nil {
		item, err := h.liveWorkloadSvc.GetPod(c.Request.Context(), name)
		if err == nil {
			c.JSON(http.StatusOK, item)
			return
		}
		if strings.Contains(err.Error(), "not found:") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	item, err := h.workloadSvc.GetPod(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *WorkloadHandler) GetPodYAML(c *gin.Context) {
	if h.adapterMode != "mock" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "pod yaml path not enabled in real-only mode"})
		return
	}
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
	if h.adapterMode != "mock" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "pod yaml write path not enabled in real-only mode"})
		return
	}
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
	query := service.PodLogQuery{
		Container:     c.Query("container"),
		Keyword:       c.Query("keyword"),
		CaseSensitive: parseBool(c.Query("caseSensitive")),
		MatchOnly:     parseBool(c.Query("matchOnly")),
		Follow:        parseBool(c.Query("follow")),
	}
	if h.adapterMode != "mock" && h.liveWorkloadSvc != nil {
		logs, err := h.liveWorkloadSvc.GetPodLogs(c.Request.Context(), name, query)
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(http.StatusOK, logs)
		return
	}

	logs, err := h.workloadSvc.GetPodLogs(name, service.PodLogQuery{
		Container:     query.Container,
		Keyword:       query.Keyword,
		CaseSensitive: query.CaseSensitive,
		MatchOnly:     query.MatchOnly,
		Follow:        query.Follow,
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
	if h.adapterMode != "mock" && h.liveWorkloadSvc != nil {
		caps, err := h.liveWorkloadSvc.GetTerminalCapabilities(c.Request.Context(), name)
		if err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, caps)
		return
	}

	caps, err := h.workloadSvc.GetTerminalCapabilities(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, caps)
}

func (h *WorkloadHandler) CreateTerminalSession(c *gin.Context) {
	name := c.Param("name")
	var req createTerminalSessionRequest
	_ = c.ShouldBindJSON(&req)

	if h.adapterMode != "mock" && h.liveWorkloadSvc != nil {
		if err := h.liveWorkloadSvc.CreateTerminalSession(c.Request.Context(), name, req.Container); err != nil {
			if strings.Contains(err.Error(), "not found:") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		session := h.terminalSessions.Create(name, req.Container, c.GetString("km_user"), c.GetString("km_role"))
		ttlSeconds := int(h.terminalSessions.TTL().Seconds())
		c.JSON(http.StatusCreated, gin.H{
			"enabled":    true,
			"sessionId":  session.ID,
			"container":  session.Container,
			"wsPath":     "/api/v1/pods/" + name + "/terminal/ws?sessionId=" + session.ID,
			"ttlSeconds": ttlSeconds,
			"expiresAt":  session.ExpiresAt.Format(time.RFC3339),
			"error":      "terminal session created",
		})
		return
	}

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
	if h.adapterMode != "mock" && h.liveWorkloadSvc != nil {
		items, err := h.liveWorkloadSvc.ListStatefulSets(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.workloadSvc.ListStatefulSets()})
}

func (h *WorkloadHandler) GetStatefulSet(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveWorkloadSvc != nil {
		item, err := h.liveWorkloadSvc.GetStatefulSet(c.Request.Context(), name)
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
	item, err := h.workloadSvc.GetStatefulSet(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *WorkloadHandler) GetStatefulSetYAML(c *gin.Context) {
	if h.adapterMode != "mock" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "statefulset yaml path not enabled in real-only mode"})
		return
	}
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
	if h.adapterMode != "mock" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "statefulset yaml write path not enabled in real-only mode"})
		return
	}
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
	if h.adapterMode != "mock" && h.liveWorkloadSvc != nil {
		items, err := h.liveWorkloadSvc.ListDaemonSets(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.workloadSvc.ListDaemonSets()})
}

func (h *WorkloadHandler) GetDaemonSet(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveWorkloadSvc != nil {
		item, err := h.liveWorkloadSvc.GetDaemonSet(c.Request.Context(), name)
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
	item, err := h.workloadSvc.GetDaemonSet(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *WorkloadHandler) GetDaemonSetYAML(c *gin.Context) {
	if h.adapterMode != "mock" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "daemonset yaml path not enabled in real-only mode"})
		return
	}
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
	if h.adapterMode != "mock" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "daemonset yaml write path not enabled in real-only mode"})
		return
	}
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
	if h.adapterMode != "mock" && h.liveWorkloadSvc != nil {
		items, err := h.liveWorkloadSvc.ListJobs(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.workloadSvc.ListJobs()})
}

func (h *WorkloadHandler) GetJob(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveWorkloadSvc != nil {
		item, err := h.liveWorkloadSvc.GetJob(c.Request.Context(), name)
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
	item, err := h.workloadSvc.GetJob(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *WorkloadHandler) GetJobYAML(c *gin.Context) {
	if h.adapterMode != "mock" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "job yaml path not enabled in real-only mode"})
		return
	}
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
	if h.adapterMode != "mock" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "job yaml write path not enabled in real-only mode"})
		return
	}
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
	if h.adapterMode != "mock" && h.liveWorkloadSvc != nil {
		items, err := h.liveWorkloadSvc.ListCronJobs(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": h.workloadSvc.ListCronJobs()})
}

func (h *WorkloadHandler) GetCronJob(c *gin.Context) {
	name := c.Param("name")
	if h.adapterMode != "mock" && h.liveWorkloadSvc != nil {
		item, err := h.liveWorkloadSvc.GetCronJob(c.Request.Context(), name)
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
	if h.adapterMode != "mock" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "cronjob yaml path not enabled in real-only mode"})
		return
	}
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
	if h.adapterMode != "mock" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "cronjob yaml write path not enabled in real-only mode"})
		return
	}
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
