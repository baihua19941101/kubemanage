package server

import (
	"kubeManage/backend/internal/handlers"
	"kubeManage/backend/internal/infra"
	"kubeManage/backend/internal/middleware"
	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(store *infra.Store) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	authSvc := service.NewAuthService()
	auditSvc := service.NewAuditService()
	r.Use(middleware.InjectRole(authSvc))

	clusterSvc := service.NewClusterService(nil)
	if store != nil {
		clusterSvc = service.NewClusterService(store.Redis)
	}
	clusterHandler := handlers.NewClusterHandler(clusterSvc)
	namespaceHandler := handlers.NewNamespaceHandler(service.NewNamespaceService())
	workloadHandler := handlers.NewWorkloadHandler(service.NewWorkloadService())
	resourceHandler := handlers.NewResourceHandler(service.NewResourceService())
	authHandler := handlers.NewAuthHandler(authSvc)
	auditHandler := handlers.NewAuditHandler(auditSvc)

	api := r.Group("/api/v1")
	{
		api.GET("/healthz", handlers.Healthz)
		api.GET("/clusters", clusterHandler.ListClusters)
		api.GET("/clusters/current", clusterHandler.GetCurrentCluster)
		api.GET("/namespaces", namespaceHandler.ListNamespaces)
		api.GET("/namespaces/:name/yaml", namespaceHandler.GetNamespaceYAML)
		api.GET("/namespaces/:name/yaml/download", namespaceHandler.DownloadNamespaceYAML)
		api.GET("/deployments", workloadHandler.ListDeployments)
		api.GET("/deployments/:name", workloadHandler.GetDeployment)
		api.GET("/deployments/:name/yaml", workloadHandler.GetDeploymentYAML)
		api.GET("/pods", workloadHandler.ListPods)
		api.GET("/pods/:name", workloadHandler.GetPod)
		api.GET("/pods/:name/yaml", workloadHandler.GetPodYAML)
		api.GET("/pods/:name/logs", workloadHandler.GetPodLogs)
		api.GET("/services", resourceHandler.ListServices)
		api.GET("/services/:name", resourceHandler.GetService)
		api.GET("/configmaps", resourceHandler.ListConfigMaps)
		api.GET("/configmaps/:name", resourceHandler.GetConfigMap)
		api.GET("/secrets", resourceHandler.ListSecrets)
		api.GET("/secrets/:name", resourceHandler.GetSecret)
		api.GET("/auth/me", authHandler.GetMe)
		api.GET("/audits", middleware.RequirePermission(authSvc, service.PermAuditRead), auditHandler.ListAudits)
	}

	write := api.Group("", middleware.WriteAudit(auditSvc))
	{
		write.POST("/clusters/switch", middleware.RequirePermission(authSvc, service.PermWorkloadWrite), clusterHandler.SwitchCluster)
		write.POST("/namespaces", middleware.RequirePermission(authSvc, service.PermNamespaceWrite), namespaceHandler.CreateNamespace)
		write.DELETE("/namespaces/:name", middleware.RequirePermission(authSvc, service.PermNamespaceWrite), namespaceHandler.DeleteNamespace)
		write.PUT("/deployments/:name/yaml", middleware.RequirePermission(authSvc, service.PermWorkloadWrite), workloadHandler.UpdateDeploymentYAML)
		write.PUT("/pods/:name/yaml", middleware.RequirePermission(authSvc, service.PermWorkloadWrite), workloadHandler.UpdatePodYAML)
	}

	return r
}
