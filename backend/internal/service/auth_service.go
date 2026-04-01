package service

type Permission string

const (
	PermNamespaceWrite Permission = "namespace:write"
	PermWorkloadWrite  Permission = "workload:write"
	PermAuditRead      Permission = "audit:read"
)

const (
	RoleViewer   = "viewer"
	RoleOperator = "operator"
	RoleAdmin    = "admin"
)

type AuthService struct {
	rolePerms      map[string]map[Permission]bool
	roleNamespaces map[string][]string
}

func NewAuthService() *AuthService {
	return &AuthService{
		rolePerms: map[string]map[Permission]bool{
			RoleViewer: {},
			RoleOperator: {
				PermNamespaceWrite: true,
				PermWorkloadWrite:  true,
			},
			RoleAdmin: {
				PermNamespaceWrite: true,
				PermWorkloadWrite:  true,
				PermAuditRead:      true,
			},
		},
		roleNamespaces: map[string][]string{
			RoleViewer:   nil,
			RoleOperator: {"dev"},
			RoleAdmin:    {"*"},
		},
	}
}

func (s *AuthService) NormalizeRole(role string) string {
	switch role {
	case RoleAdmin, RoleOperator, RoleViewer:
		return role
	default:
		return RoleViewer
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
	normalized := s.NormalizeRole(role)
	allowed := s.roleNamespaces[normalized]
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
