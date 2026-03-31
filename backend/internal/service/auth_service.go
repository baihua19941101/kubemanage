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
	rolePerms map[string]map[Permission]bool
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
