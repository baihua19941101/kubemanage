package server

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"kubeManage/backend/internal/handlers"
	"kubeManage/backend/internal/infra"
	"kubeManage/backend/internal/middleware"
	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(store *infra.Store, k8sAdapterMode string, secretKey string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.InjectRequestID())
	adapterMode := normalizeK8sAdapterMode(k8sAdapterMode)

	authSvc := service.NewAuthService()
	if store != nil && store.DB != nil {
		authSvc = service.NewAuthServiceWithStore(
			store.DB,
			strings.TrimSpace(os.Getenv("KM_AUTH_JWT_SECRET")),
			parseDurationSeconds("KM_AUTH_ACCESS_TTL_SECONDS", 3600),
			parseDurationSeconds("KM_AUTH_REFRESH_TTL_SECONDS", 604800),
		)
		_ = authSvc.EnsureDefaultAdmin(context.Background())
	}
	auditSvc := service.NewAuditService()
	r.Use(middleware.InjectRole(authSvc))

	clusterSvc := service.NewClusterService(nil)
	if store != nil {
		clusterSvc = service.NewClusterService(store.Redis)
	}
	var clusterConnectionRepo service.ClusterConnectionRepository
	var clusterConnectionAdapter service.K8sAdapter
	if store != nil && store.DB != nil {
		clusterConnectionRepo = service.NewGormClusterConnectionRepo(store.DB, secretKey)
		if adapterMode == "live" || adapterMode == "auto" {
			clusterConnectionAdapter = service.NewRealK8sAdapter()
		}
	}
	clusterConnectionSvc := service.NewClusterConnectionService(clusterConnectionRepo, clusterConnectionAdapter)
	clusterHandler := handlers.NewClusterHandler(clusterSvc, clusterConnectionSvc, adapterMode)
	clusterConnectionHandler := handlers.NewClusterConnectionHandler(clusterConnectionSvc)
	namespaceSvc := service.NewNamespaceService()
	workloadSvc := service.NewWorkloadService()
	liveWorkloadSvc := service.NewLiveWorkloadReader(clusterConnectionRepo)
	terminalSessions := service.NewTerminalSessionStore(parseTerminalSessionTTL())
	liveResourceSvc := service.NewLiveResourceReader(clusterConnectionRepo)
	namespaceHandler := handlers.NewNamespaceHandler(namespaceSvc, clusterConnectionSvc, adapterMode)
	workloadHandler := handlers.NewWorkloadHandler(workloadSvc, liveWorkloadSvc, terminalSessions, adapterMode)
	resourceHandler := handlers.NewResourceHandler(service.NewResourceService(), liveResourceSvc, adapterMode)
	authHandler := handlers.NewAuthHandler(authSvc)
	auditHandler := handlers.NewAuditHandler(auditSvc)

	api := r.Group("/api/v1")
	{
		api.GET("/healthz", handlers.Healthz)
		api.GET("/clusters", clusterHandler.ListClusters)
		api.GET("/clusters/current", clusterHandler.GetCurrentCluster)
		api.GET("/clusters/connections", clusterConnectionHandler.ListConnections)
		api.GET("/clusters/live", clusterConnectionHandler.GetLiveCluster)
		api.GET("/namespaces", namespaceHandler.ListNamespaces)
		api.GET("/namespaces/live", clusterConnectionHandler.ListLiveNamespaces)
		api.GET("/namespaces/:name/yaml", namespaceHandler.GetNamespaceYAML)
		api.GET("/namespaces/:name/yaml/download", namespaceHandler.DownloadNamespaceYAML)
		api.GET("/deployments", workloadHandler.ListDeployments)
		api.GET("/deployments/:name", workloadHandler.GetDeployment)
		api.GET("/deployments/:name/yaml", workloadHandler.GetDeploymentYAML)
		api.GET("/pods", workloadHandler.ListPods)
		api.GET("/pods/:name", workloadHandler.GetPod)
		api.GET("/pods/:name/yaml", workloadHandler.GetPodYAML)
		api.GET("/pods/:name/logs", workloadHandler.GetPodLogs)
		api.GET("/pods/:name/terminal/capabilities", workloadHandler.GetTerminalCapabilities)
		api.GET("/pods/:name/terminal/ws", workloadHandler.TerminalWebSocket)
		api.GET("/statefulsets", workloadHandler.ListStatefulSets)
		api.GET("/statefulsets/:name", workloadHandler.GetStatefulSet)
		api.GET("/statefulsets/:name/yaml", workloadHandler.GetStatefulSetYAML)
		api.GET("/daemonsets", workloadHandler.ListDaemonSets)
		api.GET("/daemonsets/:name", workloadHandler.GetDaemonSet)
		api.GET("/daemonsets/:name/yaml", workloadHandler.GetDaemonSetYAML)
		api.GET("/jobs", workloadHandler.ListJobs)
		api.GET("/jobs/:name", workloadHandler.GetJob)
		api.GET("/jobs/:name/yaml", workloadHandler.GetJobYAML)
		api.GET("/cronjobs", workloadHandler.ListCronJobs)
		api.GET("/cronjobs/:name", workloadHandler.GetCronJob)
		api.GET("/cronjobs/:name/yaml", workloadHandler.GetCronJobYAML)
		api.GET("/services", resourceHandler.ListServices)
		api.GET("/services/:name", resourceHandler.GetService)
		api.GET("/configmaps", resourceHandler.ListConfigMaps)
		api.GET("/configmaps/:name", resourceHandler.GetConfigMap)
		api.GET("/secrets", resourceHandler.ListSecrets)
		api.GET("/secrets/:name", resourceHandler.GetSecret)
		api.GET("/ingresses", resourceHandler.ListIngresses)
		api.GET("/ingresses/:name", resourceHandler.GetIngress)
		api.GET("/ingresses/:name/services", resourceHandler.ListIngressServices)
		api.GET("/hpas", resourceHandler.ListHPAs)
		api.GET("/hpas/:name", resourceHandler.GetHPA)
		api.GET("/hpas/:name/target", resourceHandler.GetHPATarget)
		api.GET("/pvs", resourceHandler.ListPVs)
		api.GET("/pvs/:name", resourceHandler.GetPV)
		api.GET("/pvcs", resourceHandler.ListPVCs)
		api.GET("/pvcs/:name", resourceHandler.GetPVC)
		api.GET("/storageclasses", resourceHandler.ListStorageClasses)
		api.GET("/storageclasses/:name", resourceHandler.GetStorageClass)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/refresh", authHandler.Refresh)
		api.GET("/auth/me", authHandler.GetMe)
		api.GET("/auth/users", middleware.RequirePermission(authSvc, service.PermUserManage), authHandler.ListUsers)
		api.GET("/audits", middleware.RequirePermission(authSvc, service.PermAuditRead), auditHandler.ListAudits)
	}

	write := api.Group("", middleware.WriteAudit(auditSvc))
	{
		write.POST("/clusters/switch", middleware.RequirePermission(authSvc, service.PermWorkloadWrite), middleware.RequireActionConfirm("switch_cluster"), clusterHandler.SwitchCluster)
		write.POST("/clusters/connections/import/kubeconfig", middleware.RequirePermission(authSvc, service.PermClusterManage), middleware.RequireActionConfirm("import_cluster_kubeconfig"), clusterConnectionHandler.ImportKubeconfig)
		write.POST("/clusters/connections/import/token", middleware.RequirePermission(authSvc, service.PermClusterManage), middleware.RequireActionConfirm("import_cluster_token"), clusterConnectionHandler.ImportToken)
		write.POST("/clusters/connections/test", middleware.RequirePermission(authSvc, service.PermClusterManage), clusterConnectionHandler.TestConnection)
		write.POST("/clusters/connections/:id/activate", middleware.RequirePermission(authSvc, service.PermClusterManage), middleware.RequireActionConfirm("activate_cluster_connection"), clusterConnectionHandler.Activate)
		write.POST("/auth/logout", authHandler.Logout)
		write.POST("/auth/users", middleware.RequirePermission(authSvc, service.PermUserManage), middleware.RequireActionConfirm("create_user"), authHandler.CreateUser)
		write.PATCH("/auth/users/:username/status", middleware.RequirePermission(authSvc, service.PermUserManage), middleware.RequireActionConfirm("update_user_status"), authHandler.UpdateUserStatus)
		write.PATCH("/auth/users/:username", middleware.RequirePermission(authSvc, service.PermUserManage), middleware.RequireActionConfirm("update_user_profile"), authHandler.UpdateUserProfile)
		write.POST("/auth/users/:username/reset-password", middleware.RequirePermission(authSvc, service.PermUserManage), middleware.RequireActionConfirm("reset_user_password"), authHandler.ResetUserPassword)
		write.POST("/namespaces", middleware.RequireScopedPermission(authSvc, service.PermNamespaceWrite, middleware.ResolvePathParamFromBodyOrJSON("name")), middleware.RequireActionConfirm("create_namespace"), namespaceHandler.CreateNamespace)
		write.DELETE("/namespaces/:name", middleware.RequireScopedPermission(authSvc, service.PermNamespaceWrite, middleware.ResolvePathParam("name")), middleware.RequireActionConfirm("delete_namespace"), namespaceHandler.DeleteNamespace)
		write.PUT("/deployments/:name/yaml", middleware.RequireScopedPermission(authSvc, service.PermWorkloadWrite, func(c *gin.Context) (string, error) {
			return workloadSvc.DeploymentNamespace(c.Param("name"))
		}), middleware.RequireActionConfirm("update_deployment_yaml"), workloadHandler.UpdateDeploymentYAML)
		write.PUT("/pods/:name/yaml", middleware.RequireScopedPermission(authSvc, service.PermWorkloadWrite, func(c *gin.Context) (string, error) {
			if adapterMode != "mock" && liveWorkloadSvc != nil {
				return liveWorkloadSvc.PodNamespace(c.Request.Context(), c.Param("name"))
			}
			return workloadSvc.PodNamespace(c.Param("name"))
		}), middleware.RequireActionConfirm("update_pod_yaml"), workloadHandler.UpdatePodYAML)
		write.POST("/pods/:name/terminal/sessions", middleware.RequireScopedPermission(authSvc, service.PermWorkloadWrite, func(c *gin.Context) (string, error) {
			if adapterMode != "mock" && liveWorkloadSvc != nil {
				return liveWorkloadSvc.PodNamespace(c.Request.Context(), c.Param("name"))
			}
			return workloadSvc.PodNamespace(c.Param("name"))
		}), middleware.RequireActionConfirm("create_terminal_session"), workloadHandler.CreateTerminalSession)
		write.PUT("/statefulsets/:name/yaml", middleware.RequireScopedPermission(authSvc, service.PermWorkloadWrite, func(c *gin.Context) (string, error) {
			return workloadSvc.StatefulSetNamespace(c.Param("name"))
		}), middleware.RequireActionConfirm("update_statefulset_yaml"), workloadHandler.UpdateStatefulSetYAML)
		write.PUT("/daemonsets/:name/yaml", middleware.RequireScopedPermission(authSvc, service.PermWorkloadWrite, func(c *gin.Context) (string, error) {
			return workloadSvc.DaemonSetNamespace(c.Param("name"))
		}), middleware.RequireActionConfirm("update_daemonset_yaml"), workloadHandler.UpdateDaemonSetYAML)
		write.PUT("/jobs/:name/yaml", middleware.RequireScopedPermission(authSvc, service.PermWorkloadWrite, func(c *gin.Context) (string, error) {
			return workloadSvc.JobNamespace(c.Param("name"))
		}), middleware.RequireActionConfirm("update_job_yaml"), workloadHandler.UpdateJobYAML)
		write.PUT("/cronjobs/:name/yaml", middleware.RequireScopedPermission(authSvc, service.PermWorkloadWrite, func(c *gin.Context) (string, error) {
			return workloadSvc.CronJobNamespace(c.Param("name"))
		}), middleware.RequireActionConfirm("update_cronjob_yaml"), workloadHandler.UpdateCronJobYAML)
	}

	return r
}

func normalizeK8sAdapterMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "live", "auto", "mock":
		return strings.ToLower(strings.TrimSpace(mode))
	default:
		return "live"
	}
}

func parseTerminalSessionTTL() time.Duration {
	const defaultTTL = 120
	value := strings.TrimSpace(os.Getenv("KM_TERMINAL_SESSION_TTL_SECONDS"))
	if value == "" {
		return time.Duration(defaultTTL) * time.Second
	}
	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return time.Duration(defaultTTL) * time.Second
	}
	return time.Duration(seconds) * time.Second
}

func parseDurationSeconds(envKey string, fallback int) time.Duration {
	value := strings.TrimSpace(os.Getenv(envKey))
	if value == "" {
		return time.Duration(fallback) * time.Second
	}
	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return time.Duration(fallback) * time.Second
	}
	return time.Duration(seconds) * time.Second
}
