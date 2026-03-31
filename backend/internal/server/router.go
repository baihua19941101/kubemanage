package server

import (
	"kubeManage/backend/internal/handlers"
	"kubeManage/backend/internal/infra"
	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(store *infra.Store) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	clusterSvc := service.NewClusterService(nil)
	if store != nil {
		clusterSvc = service.NewClusterService(store.Redis)
	}
	clusterHandler := handlers.NewClusterHandler(clusterSvc)
	namespaceHandler := handlers.NewNamespaceHandler(service.NewNamespaceService())
	workloadHandler := handlers.NewWorkloadHandler(service.NewWorkloadService())

	api := r.Group("/api/v1")
	{
		api.GET("/healthz", handlers.Healthz)
		api.GET("/clusters", clusterHandler.ListClusters)
		api.GET("/clusters/current", clusterHandler.GetCurrentCluster)
		api.POST("/clusters/switch", clusterHandler.SwitchCluster)
		api.GET("/namespaces", namespaceHandler.ListNamespaces)
		api.POST("/namespaces", namespaceHandler.CreateNamespace)
		api.DELETE("/namespaces/:name", namespaceHandler.DeleteNamespace)
		api.GET("/namespaces/:name/yaml", namespaceHandler.GetNamespaceYAML)
		api.GET("/namespaces/:name/yaml/download", namespaceHandler.DownloadNamespaceYAML)
		api.GET("/deployments", workloadHandler.ListDeployments)
		api.GET("/deployments/:name", workloadHandler.GetDeployment)
		api.GET("/deployments/:name/yaml", workloadHandler.GetDeploymentYAML)
		api.PUT("/deployments/:name/yaml", workloadHandler.UpdateDeploymentYAML)
		api.GET("/pods", workloadHandler.ListPods)
		api.GET("/pods/:name", workloadHandler.GetPod)
		api.GET("/pods/:name/yaml", workloadHandler.GetPodYAML)
		api.PUT("/pods/:name/yaml", workloadHandler.UpdatePodYAML)
		api.GET("/pods/:name/logs", workloadHandler.GetPodLogs)
	}

	return r
}
