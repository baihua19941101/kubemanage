package handlers

import (
	"net/http"
	"strconv"

	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	auditSvc *service.AuditService
}

func NewAuditHandler(auditSvc *service.AuditService) *AuditHandler {
	return &AuditHandler{auditSvc: auditSvc}
}

func (h *AuditHandler) ListAudits(c *gin.Context) {
	filter := service.AuditFilter{
		User:   c.Query("user"),
		Role:   c.Query("role"),
		Method: c.Query("method"),
		Path:   c.Query("path"),
	}
	if code, err := strconv.Atoi(c.Query("statusCode")); err == nil {
		filter.StatusCode = code
	}
	if limit, err := strconv.Atoi(c.Query("limit")); err == nil {
		filter.Limit = limit
	}
	c.JSON(http.StatusOK, gin.H{"items": h.auditSvc.List(filter)})
}
