package handlers

import (
	"net/http"
	"slices"
	"strings"

	"kubeManage/backend/internal/middleware"
	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	role := c.GetString(middleware.CtxRoleKey)
	perms := h.authSvc.Permissions(role)
	permStrings := make([]string, 0, len(perms))
	for _, p := range perms {
		permStrings = append(permStrings, string(p))
	}
	slices.Sort(permStrings)

	c.JSON(http.StatusOK, gin.H{
		"user":              c.GetString(middleware.CtxUserKey),
		"role":              role,
		"permissions":       strings.Join(permStrings, ","),
		"allowedNamespaces": h.authSvc.AllowedNamespaces(role),
	})
}
