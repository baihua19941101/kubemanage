package handlers

import (
	"net/http"
	"strconv"

	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ClusterConnectionHandler struct {
	svc *service.ClusterConnectionService
}

func NewClusterConnectionHandler(svc *service.ClusterConnectionService) *ClusterConnectionHandler {
	return &ClusterConnectionHandler{svc: svc}
}

type importKubeconfigRequest struct {
	Name              string `json:"name"`
	KubeconfigContent string `json:"kubeconfigContent"`
}

type importTokenRequest struct {
	Name          string `json:"name"`
	APIServer     string `json:"apiServer"`
	BearerToken   string `json:"bearerToken"`
	CACert        string `json:"caCert"`
	SkipTLSVerify bool   `json:"skipTlsVerify"`
}

type testConnectionRequest struct {
	Mode              string `json:"mode"`
	APIServer         string `json:"apiServer"`
	KubeconfigContent string `json:"kubeconfigContent"`
	BearerToken       string `json:"bearerToken"`
	CACert            string `json:"caCert"`
	SkipTLSVerify     bool   `json:"skipTlsVerify"`
}

func (h *ClusterConnectionHandler) ListConnections(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ClusterConnectionHandler) ImportKubeconfig(c *gin.Context) {
	var req importKubeconfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	item, err := h.svc.ImportKubeconfig(c.Request.Context(), service.ImportKubeconfigInput{Name: req.Name, KubeconfigContent: req.KubeconfigContent})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *ClusterConnectionHandler) ImportToken(c *gin.Context) {
	var req importTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	item, err := h.svc.ImportToken(c.Request.Context(), service.ImportTokenInput{Name: req.Name, APIServer: req.APIServer, BearerToken: req.BearerToken, CACert: req.CACert, SkipTLSVerify: req.SkipTLSVerify})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *ClusterConnectionHandler) TestConnection(c *gin.Context) {
	var req testConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	result, err := h.svc.TestConnection(c.Request.Context(), service.ConnectionTestInput{Mode: req.Mode, APIServer: req.APIServer, KubeconfigContent: req.KubeconfigContent, BearerToken: req.BearerToken, CACert: req.CACert, SkipTLSVerify: req.SkipTLSVerify})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ClusterConnectionHandler) Activate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Activate(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ClusterConnectionHandler) GetLiveCluster(c *gin.Context) {
	item, err := h.svc.GetLiveCluster(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ClusterConnectionHandler) ListLiveNamespaces(c *gin.Context) {
	items, err := h.svc.ListLiveNamespaces(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
