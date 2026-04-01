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
		errMsg := c.GetString("km_error")
		if errMsg == "" && len(c.Errors) > 0 {
			errMsg = c.Errors.Last().Error()
		}
		auditSvc.Append(
			c.GetString(CtxRequestIDKey),
			c.GetString(CtxUserKey),
			c.GetString(CtxRoleKey),
			c.Request.Method,
			c.Request.URL.Path,
			c.GetString("km_namespace"),
			c.Writer.Status(),
			errMsg,
		)
	}
}
