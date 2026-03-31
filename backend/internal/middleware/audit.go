package middleware

import (
	"net/http"

	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func WriteAudit(auditSvc *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Request.Method == http.MethodGet {
			return
		}
		auditSvc.Append(
			c.GetString(CtxUserKey),
			c.GetString(CtxRoleKey),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
		)
	}
}
