package handlers

import (
	"net/http"
	"slices"
	"strconv"
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
	Provider string `json:"provider"`
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

type createAuthProviderRequest struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

type updateAuthProviderStatusRequest struct {
	IsEnabled bool `json:"isEnabled"`
}

type revokeAllTokensRequest struct {
	Username string `json:"username"`
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
	pair, err := h.authSvc.Login(c.Request.Context(), req.Username, req.Password, req.Provider)
	if err != nil {
		switch err {
		case service.ErrAuthDBNotEnabled:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case service.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case service.ErrUserDisabled:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case service.ErrAuthProviderNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case service.ErrAuthProviderDisabled:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case service.ErrLDAPConfigInvalid:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case service.ErrLDAPUnavailable:
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
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

func (h *AuthHandler) ListAuthProviders(c *gin.Context) {
	items, err := h.authSvc.ListAuthProviders(c.Request.Context())
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

func (h *AuthHandler) ListPublicAuthProviders(c *gin.Context) {
	items, err := h.authSvc.ListPublicAuthProviders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *AuthHandler) ListTokenSessions(c *gin.Context) {
	username := strings.TrimSpace(c.Query("username"))
	activeOnly := false
	if raw := strings.ToLower(strings.TrimSpace(c.Query("activeOnly"))); raw == "1" || raw == "true" || raw == "yes" {
		activeOnly = true
	}
	limit := 100
	if rawLimit := strings.TrimSpace(c.Query("limit")); rawLimit != "" {
		if parsed, err := strconv.Atoi(rawLimit); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	items, err := h.authSvc.ListTokenSessions(c.Request.Context(), username, activeOnly, limit)
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

func (h *AuthHandler) RevokeTokenSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token session id"})
		return
	}
	if err := h.authSvc.RevokeTokenSession(c.Request.Context(), uint(id)); err != nil {
		switch err {
		case service.ErrAuthDBNotEnabled:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case service.ErrTokenSessionNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) RevokeAllTokens(c *gin.Context) {
	requester := strings.TrimSpace(c.GetString(middleware.CtxUserKey))
	role := strings.TrimSpace(c.GetString(middleware.CtxRoleKey))
	target := requester

	var req revokeAllTokensRequest
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		if strings.TrimSpace(req.Username) != "" {
			target = strings.TrimSpace(req.Username)
		}
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}
	if h.authSvc.NormalizeRole(role) != service.RoleAdmin && !strings.EqualFold(target, requester) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	count, err := h.authSvc.RevokeAllTokensByUser(c.Request.Context(), target)
	if err != nil {
		switch err {
		case service.ErrAuthDBNotEnabled:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"revoked": count, "username": target})
}

func (h *AuthHandler) CreateAuthProvider(c *gin.Context) {
	var req createAuthProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.authSvc.CreateAuthProvider(c.Request.Context(), req.Name, req.Type, req.Config); err != nil {
		switch err {
		case service.ErrAuthDBNotEnabled:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case service.ErrAuthProviderType, service.ErrAuthProviderName:
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

func (h *AuthHandler) UpdateAuthProviderStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider id"})
		return
	}
	var req updateAuthProviderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.authSvc.SetAuthProviderEnabled(c.Request.Context(), uint(id), req.IsEnabled); err != nil {
		switch err {
		case service.ErrAuthDBNotEnabled:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case service.ErrAuthProviderNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) SetDefaultAuthProvider(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider id"})
		return
	}
	if err := h.authSvc.SetDefaultAuthProvider(c.Request.Context(), uint(id)); err != nil {
		switch err {
		case service.ErrAuthDBNotEnabled:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case service.ErrAuthProviderNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case service.ErrAuthProviderDisabled:
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
