package middleware

import (
	"net/http"

	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	CtxRoleKey = "km_role"
	CtxUserKey = "km_user"
)

func InjectRole(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetHeader("X-User")
		if user == "" {
			user = "demo-user"
		}
		role := authSvc.NormalizeRole(c.GetHeader("X-User-Role"))
		c.Set(CtxUserKey, user)
		c.Set(CtxRoleKey, role)
		c.Next()
	}
}

func RequirePermission(authSvc *service.AuthService, perm service.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString(CtxRoleKey)
		if !authSvc.HasPermission(role, perm) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "permission denied",
				"role":  role,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
