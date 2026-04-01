package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func RequireScopedPermission(
	authSvc *service.AuthService,
	perm service.Permission,
	namespaceResolver func(*gin.Context) (string, error),
) gin.HandlerFunc {
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

		namespace, err := namespaceResolver(c)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if !authSvc.CanAccessNamespace(role, namespace) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":     "namespace access denied",
				"role":      role,
				"namespace": namespace,
			})
			c.Abort()
			return
		}

		c.Set("km_namespace", namespace)
		c.Next()
	}
}

func ResolvePathParam(param string) func(*gin.Context) (string, error) {
	return func(c *gin.Context) (string, error) {
		value := c.Param(param)
		if value == "" {
			return "", fmt.Errorf("%s is required", param)
		}
		return value, nil
	}
}

func ResolvePathParamFromBodyOrJSON(field string) func(*gin.Context) (string, error) {
	return func(c *gin.Context) (string, error) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return "", fmt.Errorf("read request body failed")
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		payload := map[string]any{}
		if err := json.Unmarshal(body, &payload); err != nil {
			return "", fmt.Errorf("invalid request body")
		}
		value, ok := payload[field].(string)
		if !ok || value == "" {
			return "", fmt.Errorf("%s is required", field)
		}
		return value, nil
	}
}
