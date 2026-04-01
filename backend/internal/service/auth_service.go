package service

import (
	"context"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"kubeManage/backend/internal/infra"

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

func (s *AuthService) Login(ctx context.Context, username, password string) (AuthTokenPair, error) {
	if s.db == nil {
		return AuthTokenPair{}, ErrAuthDBNotEnabled
	}
	var user infra.UserRecord
	if err := s.db.WithContext(ctx).Where("username = ?", strings.TrimSpace(username)).First(&user).Error; err != nil {
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
