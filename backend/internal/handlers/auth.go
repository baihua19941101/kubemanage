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

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type createUserRequest struct {
	Username          string   `json:"username"`
	Password          string   `json:"password"`
	Role              string   `json:"role"`
	AllowedNamespaces []string `json:"allowedNamespaces"`
}

type updateUserStatusRequest struct {
	IsActive bool `json:"isActive"`
}

type resetPasswordRequest struct {
	Password string `json:"password"`
}

type updateUserProfileRequest struct {
	Role              string   `json:"role"`
	AllowedNamespaces []string `json:"allowedNamespaces"`
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	pair, err := h.authSvc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		switch err {
		case service.ErrAuthDBNotEnabled:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case service.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case service.ErrUserDisabled:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, pair)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.RefreshToken) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refreshToken is required"})
		return
	}
	pair, err := h.authSvc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case service.ErrAuthDBNotEnabled:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case service.ErrRefreshTokenInvalid, service.ErrRefreshTokenRevoked:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case service.ErrUserDisabled:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, pair)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.RefreshToken) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refreshToken is required"})
		return
	}
	if err := h.authSvc.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		switch err {
		case service.ErrAuthDBNotEnabled:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case service.ErrRefreshTokenInvalid:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.authSvc.CreateUser(c.Request.Context(), req.Username, req.Password, req.Role, req.AllowedNamespaces); err != nil {
		switch err {
		case service.ErrAuthDBNotEnabled:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case service.ErrUsernameInvalid, service.ErrPasswordTooShort, service.ErrRoleNotAllowed, service.ErrReadonlyScopeRequired:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			if strings.Contains(err.Error(), "already exists") {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusCreated)
}

func (h *AuthHandler) ListUsers(c *gin.Context) {
	items, err := h.authSvc.ListUsers(c.Request.Context())
	if err != nil {
		if err == service.ErrAuthDBNotEnabled {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *AuthHandler) UpdateUserStatus(c *gin.Context) {
	username := c.Param("username")
	var req updateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.authSvc.SetUserActive(c.Request.Context(), username, req.IsActive); err != nil {
		switch err {
		case service.ErrAuthDBNotEnabled:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case service.ErrAdminDisableForbidden:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) ResetUserPassword(c *gin.Context) {
	username := c.Param("username")
	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.authSvc.ResetUserPassword(c.Request.Context(), username, req.Password); err != nil {
		switch err {
		case service.ErrAuthDBNotEnabled:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case service.ErrPasswordTooShort:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) UpdateUserProfile(c *gin.Context) {
	username := c.Param("username")
	var req updateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.authSvc.UpdateUserRoleAndNamespaces(c.Request.Context(), username, req.Role, req.AllowedNamespaces); err != nil {
		switch err {
		case service.ErrAuthDBNotEnabled:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case service.ErrRoleNotAllowed, service.ErrReadonlyScopeRequired, service.ErrAdminRoleChangeDenied:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	role := c.GetString(middleware.CtxRoleKey)
	perms := h.authSvc.Permissions(role)
	permStrings := make([]string, 0, len(perms))
	for _, p := range perms {
		permStrings = append(permStrings, string(p))
	}
	slices.Sort(permStrings)

	allowedNamespaces := c.GetStringSlice(middleware.CtxAllowedNamespacesKey)
	if len(allowedNamespaces) == 0 {
		allowedNamespaces = h.authSvc.AllowedNamespaces(role)
	}

	c.JSON(http.StatusOK, gin.H{
		"user":              c.GetString(middleware.CtxUserKey),
		"role":              role,
		"permissions":       strings.Join(permStrings, ","),
		"allowedNamespaces": allowedNamespaces,
	})
}
