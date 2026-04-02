package service

import (
	"context"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"kubeManage/backend/internal/infra"

	ldap "github.com/go-ldap/ldap/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Permission string

const (
	PermClusterManage  Permission = "cluster:manage"
	PermNamespaceWrite Permission = "namespace:write"
	PermWorkloadWrite  Permission = "workload:write"
	PermAuditRead      Permission = "audit:read"
	PermUserManage     Permission = "user:manage"
)

const (
	RoleViewer       = "viewer"
	RoleOperator     = "operator"
	RoleAdmin        = "admin"
	RoleStandardUser = "standard-user"
	RoleReadonly     = "readonly"
)

var (
	ErrAuthDBNotEnabled      = errors.New("auth database not enabled")
	ErrInvalidCredentials    = errors.New("invalid username or password")
	ErrUserDisabled          = errors.New("user is disabled")
	ErrRefreshTokenInvalid   = errors.New("refresh token invalid")
	ErrRefreshTokenRevoked   = errors.New("refresh token revoked")
	ErrRoleNotAllowed        = errors.New("role not allowed")
	ErrPasswordTooShort      = errors.New("password is too short (min 6)")
	ErrUsernameInvalid       = errors.New("username must be 3-64 chars and contain letters, digits, _ or -")
	ErrReadonlyScopeRequired = errors.New("readonly requires at least one allowed namespace")
	ErrUserNotFound          = errors.New("user not found")
	ErrAdminDisableForbidden = errors.New("admin user cannot be disabled")
	ErrAdminRoleChangeDenied = errors.New("admin user role cannot be changed")
	ErrAuthProviderNotFound  = errors.New("auth provider not found")
	ErrAuthProviderDisabled  = errors.New("auth provider is disabled")
	ErrAuthProviderType      = errors.New("auth provider type not supported")
	ErrAuthProviderName      = errors.New("auth provider name is required")
	ErrLDAPConfigInvalid     = errors.New("ldap provider config invalid")
	ErrLDAPUnavailable       = errors.New("ldap provider unavailable")
	ErrTokenSessionNotFound  = errors.New("token session not found")
)

type AuthIdentity struct {
	User              string
	Role              string
	AllowedNamespaces []string
}

type AuthTokenPair struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	TokenType    string   `json:"tokenType"`
	ExpiresIn    int64    `json:"expiresIn"`
	User         string   `json:"user"`
	Role         string   `json:"role"`
	Namespaces   []string `json:"allowedNamespaces"`
}

type UserInfo struct {
	ID                uint      `json:"id"`
	Username          string    `json:"username"`
	Role              string    `json:"role"`
	AllowedNamespaces []string  `json:"allowedNamespaces"`
	IsActive          bool      `json:"isActive"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type AuthProviderInfo struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	IsEnabled bool              `json:"isEnabled"`
	IsDefault bool              `json:"isDefault"`
	Config    map[string]string `json:"config"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

type PublicAuthProvider struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsDefault bool   `json:"isDefault"`
}

type TokenSessionInfo struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"userId"`
	Username  string     `json:"username"`
	ExpiresAt time.Time  `json:"expiresAt"`
	RevokedAt *time.Time `json:"revokedAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	Status    string     `json:"status"`
	IsActive  bool       `json:"isActive"`
}

type AuthService struct {
	rolePerms      map[string]map[Permission]bool
	roleNamespaces map[string][]string

	db         *gorm.DB
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService() *AuthService {
	return newAuthService(nil, "km-dev-jwt-secret", time.Hour, 7*24*time.Hour)
}

func NewAuthServiceWithStore(db *gorm.DB, jwtSecret string, accessTTL, refreshTTL time.Duration) *AuthService {
	return newAuthService(db, jwtSecret, accessTTL, refreshTTL)
}

func newAuthService(db *gorm.DB, jwtSecret string, accessTTL, refreshTTL time.Duration) *AuthService {
	if strings.TrimSpace(jwtSecret) == "" {
		jwtSecret = "km-dev-jwt-secret"
	}
	if accessTTL <= 0 {
		accessTTL = time.Hour
	}
	if refreshTTL <= 0 {
		refreshTTL = 7 * 24 * time.Hour
	}
	return &AuthService{
		rolePerms: map[string]map[Permission]bool{
			RoleReadonly: {},
			RoleStandardUser: {
				PermNamespaceWrite: true,
				PermWorkloadWrite:  true,
			},
			RoleAdmin: {
				PermClusterManage:  true,
				PermNamespaceWrite: true,
				PermWorkloadWrite:  true,
				PermAuditRead:      true,
				PermUserManage:     true,
			},
		},
		roleNamespaces: map[string][]string{
			RoleReadonly:     {},
			RoleStandardUser: {"dev"},
			RoleAdmin:        {"*"},
		},
		db:         db,
		jwtSecret:  []byte(jwtSecret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *AuthService) EnsureDefaultAdmin(ctx context.Context) error {
	if s.db == nil {
		return nil
	}
	var count int64
	if err := s.db.WithContext(ctx).Model(&infra.UserRecord{}).Where("username = ?", "admin").Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Create(&infra.UserRecord{
		Username:          "admin",
		PasswordHash:      string(hash),
		Role:              RoleAdmin,
		AllowedNamespaces: "*",
		IsActive:          true,
	}).Error
}

func (s *AuthService) EnsureDefaultProviders(ctx context.Context) error {
	if s.db == nil {
		return nil
	}
	var localCount int64
	if err := s.db.WithContext(ctx).Model(&infra.AuthProviderRecord{}).Where("name = ?", "local").Count(&localCount).Error; err != nil {
		return err
	}
	if localCount == 0 {
		if err := s.db.WithContext(ctx).Create(&infra.AuthProviderRecord{
			Name:      "local",
			Type:      "local",
			IsEnabled: true,
			IsDefault: true,
			Config:    "{}",
		}).Error; err != nil {
			return err
		}
	}
	var defaultCount int64
	if err := s.db.WithContext(ctx).Model(&infra.AuthProviderRecord{}).Where("is_default = ?", true).Count(&defaultCount).Error; err != nil {
		return err
	}
	if defaultCount == 0 {
		return s.db.WithContext(ctx).
			Model(&infra.AuthProviderRecord{}).
			Where("name = ?", "local").
			Updates(map[string]any{"is_default": true, "updated_at": time.Now()}).Error
	}
	return nil
}

func (s *AuthService) NormalizeRole(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case RoleAdmin:
		return RoleAdmin
	case RoleStandardUser, RoleOperator:
		return RoleStandardUser
	case RoleReadonly, RoleViewer:
		return RoleReadonly
	default:
		return RoleReadonly
	}
}

func (s *AuthService) HasPermission(role string, perm Permission) bool {
	normalized := s.NormalizeRole(role)
	return s.rolePerms[normalized][perm]
}

func (s *AuthService) Permissions(role string) []Permission {
	normalized := s.NormalizeRole(role)
	perms := make([]Permission, 0, len(s.rolePerms[normalized]))
	for p := range s.rolePerms[normalized] {
		perms = append(perms, p)
	}
	return perms
}

func (s *AuthService) AllowedNamespaces(role string) []string {
	normalized := s.NormalizeRole(role)
	namespaces := s.roleNamespaces[normalized]
	out := make([]string, len(namespaces))
	copy(out, namespaces)
	return out
}

func (s *AuthService) CanAccessNamespace(role, namespace string) bool {
	return canAccessWithAllowedNamespaces(s.AllowedNamespaces(role), namespace)
}

func CanAccessNamespace(role, namespace string, allowed []string) bool {
	normalized := strings.ToLower(strings.TrimSpace(role))
	if normalized == RoleAdmin {
		return true
	}
	return canAccessWithAllowedNamespaces(allowed, namespace)
}

func canAccessWithAllowedNamespaces(allowed []string, namespace string) bool {
	if len(allowed) == 0 {
		return false
	}
	for _, item := range allowed {
		if item == "*" || item == namespace {
			return true
		}
	}
	return false
}

func (s *AuthService) Login(ctx context.Context, username, password, provider string) (AuthTokenPair, error) {
	if s.db == nil {
		return AuthTokenPair{}, ErrAuthDBNotEnabled
	}
	providerRecord, err := s.resolveProvider(ctx, provider)
	if err != nil {
		return AuthTokenPair{}, err
	}
	providerType := normalizeProviderType(providerRecord.Type)
	username = strings.TrimSpace(username)
	if providerType == "ldap" {
		if err := s.authenticateLDAP(ctx, providerRecord, username, password); err != nil {
			return AuthTokenPair{}, err
		}
		var user infra.UserRecord
		if err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return AuthTokenPair{}, ErrInvalidCredentials
			}
			return AuthTokenPair{}, err
		}
		if !user.IsActive {
			return AuthTokenPair{}, ErrUserDisabled
		}
		return s.issueTokenPair(ctx, user)
	}
	if providerType != "local" {
		return AuthTokenPair{}, ErrAuthProviderType
	}
	var user infra.UserRecord
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AuthTokenPair{}, ErrInvalidCredentials
		}
		return AuthTokenPair{}, err
	}
	if !user.IsActive {
		return AuthTokenPair{}, ErrUserDisabled
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return AuthTokenPair{}, ErrInvalidCredentials
	}
	return s.issueTokenPair(ctx, user)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (AuthTokenPair, error) {
	if s.db == nil {
		return AuthTokenPair{}, ErrAuthDBNotEnabled
	}
	claims, err := s.parseToken(refreshToken, "refresh")
	if err != nil {
		return AuthTokenPair{}, ErrRefreshTokenInvalid
	}
	tokenHash := tokenHash(refreshToken)
	var tokenRecord infra.RefreshTokenRecord
	if err := s.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&tokenRecord).Error; err != nil {
		return AuthTokenPair{}, ErrRefreshTokenInvalid
	}
	if tokenRecord.RevokedAt != nil || time.Now().After(tokenRecord.ExpiresAt) {
		return AuthTokenPair{}, ErrRefreshTokenRevoked
	}
	userID, ok := claims["uid"].(float64)
	if !ok || uint(userID) != tokenRecord.UserID {
		return AuthTokenPair{}, ErrRefreshTokenInvalid
	}

	var user infra.UserRecord
	if err := s.db.WithContext(ctx).First(&user, tokenRecord.UserID).Error; err != nil {
		return AuthTokenPair{}, ErrRefreshTokenInvalid
	}
	if !user.IsActive {
		return AuthTokenPair{}, ErrUserDisabled
	}

	now := time.Now()
	tokenRecord.RevokedAt = &now
	if err := s.db.WithContext(ctx).Save(&tokenRecord).Error; err != nil {
		return AuthTokenPair{}, err
	}

	return s.issueTokenPair(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if s.db == nil {
		return ErrAuthDBNotEnabled
	}
	hash := tokenHash(refreshToken)
	var tokenRecord infra.RefreshTokenRecord
	if err := s.db.WithContext(ctx).Where("token_hash = ?", hash).First(&tokenRecord).Error; err != nil {
		return ErrRefreshTokenInvalid
	}
	now := time.Now()
	tokenRecord.RevokedAt = &now
	return s.db.WithContext(ctx).Save(&tokenRecord).Error
}

func (s *AuthService) CreateUser(ctx context.Context, username, password, role string, allowedNamespaces []string) error {
	if s.db == nil {
		return ErrAuthDBNotEnabled
	}
	username = strings.TrimSpace(username)
	if !validUsername(username) {
		return ErrUsernameInvalid
	}
	if len(password) < 6 {
		return ErrPasswordTooShort
	}
	normalizedRole := s.NormalizeRole(role)
	if normalizedRole != RoleAdmin && normalizedRole != RoleStandardUser && normalizedRole != RoleReadonly {
		return ErrRoleNotAllowed
	}

	normalizedAllowed := normalizeAllowedNamespaces(allowedNamespaces)
	if normalizedRole == RoleAdmin {
		normalizedAllowed = []string{"*"}
	}
	if normalizedRole == RoleReadonly && len(normalizedAllowed) == 0 {
		return ErrReadonlyScopeRequired
	}
	if normalizedRole == RoleStandardUser && len(normalizedAllowed) == 0 {
		normalizedAllowed = []string{"dev"}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	record := infra.UserRecord{
		Username:          username,
		PasswordHash:      string(hash),
		Role:              normalizedRole,
		AllowedNamespaces: strings.Join(normalizedAllowed, ","),
		IsActive:          true,
	}
	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") || strings.Contains(strings.ToLower(err.Error()), "unique") {
			return fmt.Errorf("username already exists")
		}
		return err
	}
	return nil
}

func (s *AuthService) ListUsers(ctx context.Context) ([]UserInfo, error) {
	if s.db == nil {
		return nil, ErrAuthDBNotEnabled
	}
	var records []infra.UserRecord
	if err := s.db.WithContext(ctx).Order("id asc").Find(&records).Error; err != nil {
		return nil, err
	}
	items := make([]UserInfo, 0, len(records))
	for _, item := range records {
		items = append(items, UserInfo{
			ID:                item.ID,
			Username:          item.Username,
			Role:              s.NormalizeRole(item.Role),
			AllowedNamespaces: normalizeAllowedNamespaces(strings.Split(item.AllowedNamespaces, ",")),
			IsActive:          item.IsActive,
			CreatedAt:         item.CreatedAt,
			UpdatedAt:         item.UpdatedAt,
		})
	}
	return items, nil
}

func (s *AuthService) SetUserActive(ctx context.Context, username string, active bool) error {
	if s.db == nil {
		return ErrAuthDBNotEnabled
	}
	username = strings.TrimSpace(username)
	var record infra.UserRecord
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	if strings.EqualFold(record.Username, "admin") && !active {
		return ErrAdminDisableForbidden
	}
	record.IsActive = active
	return s.db.WithContext(ctx).Save(&record).Error
}

func (s *AuthService) ResetUserPassword(ctx context.Context, username, newPassword string) error {
	if s.db == nil {
		return ErrAuthDBNotEnabled
	}
	if len(newPassword) < 6 {
		return ErrPasswordTooShort
	}
	username = strings.TrimSpace(username)
	var record infra.UserRecord
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	record.PasswordHash = string(hash)
	return s.db.WithContext(ctx).Save(&record).Error
}

func (s *AuthService) UpdateUserRoleAndNamespaces(ctx context.Context, username, role string, allowedNamespaces []string) error {
	if s.db == nil {
		return ErrAuthDBNotEnabled
	}
	username = strings.TrimSpace(username)
	var record infra.UserRecord
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	normalizedRole := s.NormalizeRole(role)
	if normalizedRole != RoleAdmin && normalizedRole != RoleStandardUser && normalizedRole != RoleReadonly {
		return ErrRoleNotAllowed
	}

	if strings.EqualFold(record.Username, "admin") && normalizedRole != RoleAdmin {
		return ErrAdminRoleChangeDenied
	}

	normalizedAllowed := normalizeAllowedNamespaces(allowedNamespaces)
	if normalizedRole == RoleAdmin {
		normalizedAllowed = []string{"*"}
	}
	if normalizedRole == RoleReadonly && len(normalizedAllowed) == 0 {
		return ErrReadonlyScopeRequired
	}
	if normalizedRole == RoleStandardUser && len(normalizedAllowed) == 0 {
		normalizedAllowed = []string{"dev"}
	}

	record.Role = normalizedRole
	record.AllowedNamespaces = strings.Join(normalizedAllowed, ",")
	return s.db.WithContext(ctx).Save(&record).Error
}

func (s *AuthService) ListAuthProviders(ctx context.Context) ([]AuthProviderInfo, error) {
	if s.db == nil {
		return nil, ErrAuthDBNotEnabled
	}
	var records []infra.AuthProviderRecord
	if err := s.db.WithContext(ctx).Order("id asc").Find(&records).Error; err != nil {
		return nil, err
	}
	items := make([]AuthProviderInfo, 0, len(records))
	for _, record := range records {
		items = append(items, AuthProviderInfo{
			ID:        record.ID,
			Name:      record.Name,
			Type:      record.Type,
			IsEnabled: record.IsEnabled,
			IsDefault: record.IsDefault,
			Config:    decodeProviderConfig(record.Config),
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
		})
	}
	return items, nil
}

func (s *AuthService) ListPublicAuthProviders(ctx context.Context) ([]PublicAuthProvider, error) {
	if s.db == nil {
		return []PublicAuthProvider{{Name: "local", Type: "local", IsDefault: true}}, nil
	}
	var records []infra.AuthProviderRecord
	if err := s.db.WithContext(ctx).Where("is_enabled = ?", true).Order("id asc").Find(&records).Error; err != nil {
		return nil, err
	}
	items := make([]PublicAuthProvider, 0, len(records))
	for _, record := range records {
		items = append(items, PublicAuthProvider{
			Name:      record.Name,
			Type:      record.Type,
			IsDefault: record.IsDefault,
		})
	}
	if len(items) == 0 {
		items = append(items, PublicAuthProvider{Name: "local", Type: "local", IsDefault: true})
	}
	return items, nil
}

func (s *AuthService) CreateAuthProvider(ctx context.Context, name, providerType string, config map[string]string) error {
	if s.db == nil {
		return ErrAuthDBNotEnabled
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrAuthProviderName
	}
	providerType = normalizeProviderType(providerType)
	if providerType != "local" && providerType != "ldap" {
		return ErrAuthProviderType
	}
	record := infra.AuthProviderRecord{
		Name:      name,
		Type:      providerType,
		IsEnabled: true,
		IsDefault: false,
		Config:    encodeProviderConfig(config),
	}
	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") || strings.Contains(strings.ToLower(err.Error()), "unique") {
			return fmt.Errorf("auth provider already exists")
		}
		return err
	}
	return nil
}

func (s *AuthService) SetAuthProviderEnabled(ctx context.Context, id uint, enabled bool) error {
	if s.db == nil {
		return ErrAuthDBNotEnabled
	}
	var record infra.AuthProviderRecord
	if err := s.db.WithContext(ctx).First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAuthProviderNotFound
		}
		return err
	}
	record.IsEnabled = enabled
	return s.db.WithContext(ctx).Save(&record).Error
}

func (s *AuthService) SetDefaultAuthProvider(ctx context.Context, id uint) error {
	if s.db == nil {
		return ErrAuthDBNotEnabled
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var record infra.AuthProviderRecord
		if err := tx.First(&record, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrAuthProviderNotFound
			}
			return err
		}
		if !record.IsEnabled {
			return ErrAuthProviderDisabled
		}
		if err := tx.Model(&infra.AuthProviderRecord{}).Where("is_default = ?", true).Updates(map[string]any{"is_default": false, "updated_at": time.Now()}).Error; err != nil {
			return err
		}
		return tx.Model(&infra.AuthProviderRecord{}).Where("id = ?", id).Updates(map[string]any{"is_default": true, "updated_at": time.Now()}).Error
	})
}

func (s *AuthService) ListTokenSessions(ctx context.Context, username string, activeOnly bool, limit int) ([]TokenSessionInfo, error) {
	if s.db == nil {
		return nil, ErrAuthDBNotEnabled
	}
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	now := time.Now()

	type tokenSessionRow struct {
		ID        uint
		UserID    uint
		Username  string
		ExpiresAt time.Time
		RevokedAt *time.Time
		CreatedAt time.Time
	}
	rows := make([]tokenSessionRow, 0, limit)
	query := s.db.WithContext(ctx).
		Table("refresh_tokens").
		Select("refresh_tokens.id, refresh_tokens.user_id, users.username, refresh_tokens.expires_at, refresh_tokens.revoked_at, refresh_tokens.created_at").
		Joins("join users on users.id = refresh_tokens.user_id")
	if strings.TrimSpace(username) != "" {
		query = query.Where("users.username = ?", strings.TrimSpace(username))
	}
	if activeOnly {
		query = query.Where("refresh_tokens.revoked_at is null and refresh_tokens.expires_at > ?", now)
	}
	if err := query.Order("refresh_tokens.id desc").Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]TokenSessionInfo, 0, len(rows))
	for _, row := range rows {
		status := "active"
		active := true
		if row.RevokedAt != nil {
			status = "revoked"
			active = false
		} else if now.After(row.ExpiresAt) {
			status = "expired"
			active = false
		}
		items = append(items, TokenSessionInfo{
			ID:        row.ID,
			UserID:    row.UserID,
			Username:  row.Username,
			ExpiresAt: row.ExpiresAt,
			RevokedAt: row.RevokedAt,
			CreatedAt: row.CreatedAt,
			Status:    status,
			IsActive:  active,
		})
	}
	return items, nil
}

func (s *AuthService) RevokeTokenSession(ctx context.Context, id uint) error {
	if s.db == nil {
		return ErrAuthDBNotEnabled
	}
	var record infra.RefreshTokenRecord
	if err := s.db.WithContext(ctx).First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTokenSessionNotFound
		}
		return err
	}
	if record.RevokedAt != nil {
		return nil
	}
	now := time.Now()
	record.RevokedAt = &now
	return s.db.WithContext(ctx).Save(&record).Error
}

func (s *AuthService) RevokeAllTokensByUser(ctx context.Context, username string) (int64, error) {
	if s.db == nil {
		return 0, ErrAuthDBNotEnabled
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return 0, ErrUserNotFound
	}
	var user infra.UserRecord
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrUserNotFound
		}
		return 0, err
	}
	now := time.Now()
	result := s.db.WithContext(ctx).
		Model(&infra.RefreshTokenRecord{}).
		Where("user_id = ? and revoked_at is null", user.ID).
		Updates(map[string]any{"revoked_at": &now, "updated_at": now})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (s *AuthService) ParseAccessToken(accessToken string) (*AuthIdentity, error) {
	claims, err := s.parseToken(accessToken, "access")
	if err != nil {
		return nil, err
	}
	user, _ := claims["usr"].(string)
	role, _ := claims["role"].(string)
	ans, _ := claims["ans"].(string)
	return &AuthIdentity{
		User:              user,
		Role:              s.NormalizeRole(role),
		AllowedNamespaces: normalizeAllowedNamespaces(strings.Split(ans, ",")),
	}, nil
}

func (s *AuthService) issueTokenPair(ctx context.Context, user infra.UserRecord) (AuthTokenPair, error) {
	now := time.Now()
	allowed := normalizeAllowedNamespaces(strings.Split(user.AllowedNamespaces, ","))
	accessClaims := jwt.MapClaims{
		"typ":  "access",
		"uid":  user.ID,
		"usr":  user.Username,
		"role": s.NormalizeRole(user.Role),
		"ans":  strings.Join(allowed, ","),
		"iat":  now.Unix(),
		"exp":  now.Add(s.accessTTL).Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.jwtSecret)
	if err != nil {
		return AuthTokenPair{}, err
	}

	refreshClaims := jwt.MapClaims{
		"typ": "refresh",
		"uid": user.ID,
		"jti": randomTokenID(12),
		"iat": now.Unix(),
		"exp": now.Add(s.refreshTTL).Unix(),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(s.jwtSecret)
	if err != nil {
		return AuthTokenPair{}, err
	}

	if s.db != nil {
		record := infra.RefreshTokenRecord{
			UserID:    user.ID,
			TokenHash: tokenHash(refreshToken),
			ExpiresAt: now.Add(s.refreshTTL),
		}
		if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
			return AuthTokenPair{}, err
		}
	}

	return AuthTokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessTTL.Seconds()),
		User:         user.Username,
		Role:         s.NormalizeRole(user.Role),
		Namespaces:   allowed,
	}, nil
}

func (s *AuthService) parseToken(rawToken string, expectType string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	typ, _ := claims["typ"].(string)
	if typ != expectType {
		return nil, fmt.Errorf("unexpected token type")
	}
	return claims, nil
}

func tokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func randomTokenID(size int) string {
	if size <= 0 {
		size = 12
	}
	buf := make([]byte, size)
	if _, err := crand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

func validUsername(username string) bool {
	if len(username) < 3 || len(username) > 64 {
		return false
	}
	for _, ch := range username {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '-' {
			continue
		}
		return false
	}
	return true
}

func normalizeAllowedNamespaces(items []string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		v := strings.TrimSpace(item)
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	slices.Sort(out)
	out = slices.Compact(out)
	return out
}

func normalizeProviderType(providerType string) string {
	return strings.ToLower(strings.TrimSpace(providerType))
}

func (s *AuthService) resolveProvider(ctx context.Context, provider string) (*infra.AuthProviderRecord, error) {
	p := strings.TrimSpace(provider)
	if p == "" {
		var defaultRecord infra.AuthProviderRecord
		if err := s.db.WithContext(ctx).Where("is_default = ?", true).First(&defaultRecord).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &infra.AuthProviderRecord{
					Name:      "local",
					Type:      "local",
					IsEnabled: true,
					IsDefault: true,
				}, nil
			}
			return nil, err
		}
		if !defaultRecord.IsEnabled {
			return nil, ErrAuthProviderDisabled
		}
		return &defaultRecord, nil
	}

	var record infra.AuthProviderRecord
	if err := s.db.WithContext(ctx).Where("name = ?", p).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAuthProviderNotFound
		}
		return nil, err
	}
	if !record.IsEnabled {
		return nil, ErrAuthProviderDisabled
	}
	return &record, nil
}

func (s *AuthService) authenticateLDAP(ctx context.Context, provider *infra.AuthProviderRecord, username, password string) error {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return ErrInvalidCredentials
	}
	cfg := decodeProviderConfig(provider.Config)
	rawURL := strings.TrimSpace(cfg["url"])
	baseDN := strings.TrimSpace(cfg["baseDN"])
	if rawURL == "" || baseDN == "" {
		return ErrLDAPConfigInvalid
	}

	loginAttr := strings.TrimSpace(cfg["loginAttr"])
	if loginAttr == "" {
		loginAttr = "uid"
	}

	filter := strings.TrimSpace(cfg["userFilter"])
	escapedUser := ldap.EscapeFilter(username)
	if filter == "" {
		filter = fmt.Sprintf("(%s=%s)", loginAttr, escapedUser)
	} else {
		filter = strings.ReplaceAll(filter, "{{username}}", escapedUser)
		filter = strings.ReplaceAll(filter, "{{loginAttr}}", loginAttr)
	}

	timeoutSeconds := 5
	if rawTimeout := strings.TrimSpace(cfg["timeoutSeconds"]); rawTimeout != "" {
		if parsed, err := strconv.Atoi(rawTimeout); err == nil && parsed > 0 {
			timeoutSeconds = parsed
		}
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return ErrLDAPConfigInvalid
	}
	timeout := time.Duration(timeoutSeconds) * time.Second
	dialer := &net.Dialer{Timeout: timeout}
	conn, err := ldap.DialURL(rawURL, ldap.DialWithDialer(dialer))
	if err != nil {
		return ErrLDAPUnavailable
	}
	defer conn.Close()
	conn.SetTimeout(timeout)

	bindDN := strings.TrimSpace(cfg["bindDN"])
	bindPassword := strings.TrimSpace(cfg["bindPassword"])
	if bindDN != "" {
		if err := conn.Bind(bindDN, bindPassword); err != nil {
			return ErrLDAPUnavailable
		}
	}

	searchReq := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 2, int(timeout.Seconds()), false,
		filter,
		[]string{"dn"},
		nil,
	)
	searchResult, err := conn.Search(searchReq)
	if err != nil {
		return ErrLDAPUnavailable
	}
	if len(searchResult.Entries) != 1 {
		return ErrInvalidCredentials
	}

	userDN := searchResult.Entries[0].DN
	if strings.TrimSpace(userDN) == "" {
		return ErrInvalidCredentials
	}
	if err := conn.Bind(userDN, password); err != nil {
		return ErrInvalidCredentials
	}
	return nil
}

func encodeProviderConfig(cfg map[string]string) string {
	if len(cfg) == 0 {
		return "{}"
	}
	parts := make([]string, 0, len(cfg))
	for k, v := range cfg {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%s", k, strings.TrimSpace(v)))
	}
	slices.Sort(parts)
	return strings.Join(parts, ";")
}

func decodeProviderConfig(raw string) map[string]string {
	out := map[string]string{}
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "{}" {
		return out
	}
	items := strings.Split(raw, ";")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		pair := strings.SplitN(item, "=", 2)
		if len(pair) != 2 {
			continue
		}
		key := strings.TrimSpace(pair[0])
		value := strings.TrimSpace(pair[1])
		if key == "" {
			continue
		}
		out[key] = value
	}
	return out
}
