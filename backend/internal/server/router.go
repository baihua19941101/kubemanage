package server

import (
	"strings"

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
	liveResourceSvc := service.NewLiveResourceReader(clusterConnectionRepo)
	namespaceHandler := handlers.NewNamespaceHandler(namespaceSvc, clusterConnectionSvc, adapterMode)
	workloadHandler := handlers.NewWorkloadHandler(workloadSvc, liveWorkloadSvc, adapterMode)
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
		api.GET("/auth/me", authHandler.GetMe)
		api.GET("/audits", middleware.RequirePermission(authSvc, service.PermAuditRead), auditHandler.ListAudits)
	}

	write := api.Group("", middleware.WriteAudit(auditSvc))
	{
		write.POST("/clusters/switch", middleware.RequirePermission(authSvc, service.PermWorkloadWrite), middleware.RequireActionConfirm("switch_cluster"), clusterHandler.SwitchCluster)
		write.POST("/clusters/connections/import/kubeconfig", middleware.RequirePermission(authSvc, service.PermClusterManage), middleware.RequireActionConfirm("import_cluster_kubeconfig"), clusterConnectionHandler.ImportKubeconfig)
		write.POST("/clusters/connections/import/token", middleware.RequirePermission(authSvc, service.PermClusterManage), middleware.RequireActionConfirm("import_cluster_token"), clusterConnectionHandler.ImportToken)
		write.POST("/clusters/connections/test", middleware.RequirePermission(authSvc, service.PermClusterManage), clusterConnectionHandler.TestConnection)
		write.POST("/clusters/connections/:id/activate", middleware.RequirePermission(authSvc, service.PermClusterManage), middleware.RequireActionConfirm("activate_cluster_connection"), clusterConnectionHandler.Activate)
		write.POST("/namespaces", middleware.RequireScopedPermission(authSvc, service.PermNamespaceWrite, middleware.ResolvePathParamFromBodyOrJSON("name")), middleware.RequireActionConfirm("create_namespace"), namespaceHandler.CreateNamespace)
		write.DELETE("/namespaces/:name", middleware.RequireScopedPermission(authSvc, service.PermNamespaceWrite, middleware.ResolvePathParam("name")), middleware.RequireActionConfirm("delete_namespace"), namespaceHandler.DeleteNamespace)
		write.PUT("/deployments/:name/yaml", middleware.RequireScopedPermission(authSvc, service.PermWorkloadWrite, func(c *gin.Context) (string, error) {
			return workloadSvc.DeploymentNamespace(c.Param("name"))
		}), middleware.RequireActionConfirm("update_deployment_yaml"), workloadHandler.UpdateDeploymentYAML)
		write.PUT("/pods/:name/yaml", middleware.RequireScopedPermission(authSvc, service.PermWorkloadWrite, func(c *gin.Context) (string, error) {
			return workloadSvc.PodNamespace(c.Param("name"))
		}), middleware.RequireActionConfirm("update_pod_yaml"), workloadHandler.UpdatePodYAML)
		write.POST("/pods/:name/terminal/sessions", middleware.RequireScopedPermission(authSvc, service.PermWorkloadWrite, func(c *gin.Context) (string, error) {
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
