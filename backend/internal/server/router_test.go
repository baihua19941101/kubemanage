package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func requestWithRole(method, path string, body string, role string) *http.Request {
	var req *http.Request
	if body != "" {
		req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}
	req.Header.Set("X-User-Role", role)
	req.Header.Set("X-User", "tester")
	if method != http.MethodGet {
		req.Header.Set("X-Action-Confirm", "CONFIRM")
	}
	return req
}

func TestHealthz(t *testing.T) {
	r := NewRouter(nil, "mock", "")
	req, err := http.NewRequest(http.MethodGet, "/api/v1/healthz", nil)
	if err != nil {
		t.Fatalf("build request failed: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestClusters(t *testing.T) {
	r := NewRouter(nil, "mock", "")
	req, err := http.NewRequest(http.MethodGet, "/api/v1/clusters", nil)
	if err != nil {
		t.Fatalf("build request failed: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	importReq := requestWithRole(http.MethodPost, "/api/v1/clusters/connections/import/token", `{"name":"dev-live","apiServer":"https://k8s.example.local","bearerToken":"token-123","caCert":"ca","skipTlsVerify":true}`, "admin")
	importW := httptest.NewRecorder()
	r.ServeHTTP(importW, importReq)
	if importW.Code != http.StatusCreated {
		t.Fatalf("import token cluster connection failed: %d body=%s", importW.Code, importW.Body.String())
	}

	listConnReq := requestWithRole(http.MethodGet, "/api/v1/clusters/connections", "", "admin")
	listConnW := httptest.NewRecorder()
	r.ServeHTTP(listConnW, listConnReq)
	if listConnW.Code != http.StatusOK {
		t.Fatalf("list cluster connections failed: %d body=%s", listConnW.Code, listConnW.Body.String())
	}

	testConnReq := requestWithRole(http.MethodPost, "/api/v1/clusters/connections/test", `{"mode":"token","apiServer":"https://k8s.example.local","bearerToken":"token-123","caCert":"ca","skipTlsVerify":true}`, "admin")
	testConnW := httptest.NewRecorder()
	r.ServeHTTP(testConnW, testConnReq)
	if testConnW.Code != http.StatusOK {
		t.Fatalf("test cluster connection failed: %d body=%s", testConnW.Code, testConnW.Body.String())
	}

	activateReq := requestWithRole(http.MethodPost, "/api/v1/clusters/connections/1/activate", "", "admin")
	activateW := httptest.NewRecorder()
	r.ServeHTTP(activateW, activateReq)
	if activateW.Code != http.StatusNoContent {
		t.Fatalf("activate cluster connection failed: %d body=%s", activateW.Code, activateW.Body.String())
	}

	liveClusterReq := requestWithRole(http.MethodGet, "/api/v1/clusters/live", "", "viewer")
	liveClusterW := httptest.NewRecorder()
	r.ServeHTTP(liveClusterW, liveClusterReq)
	if liveClusterW.Code != http.StatusOK {
		t.Fatalf("get live cluster failed: %d body=%s", liveClusterW.Code, liveClusterW.Body.String())
	}

	liveNamespacesReq := requestWithRole(http.MethodGet, "/api/v1/namespaces/live", "", "viewer")
	liveNamespacesW := httptest.NewRecorder()
	r.ServeHTTP(liveNamespacesW, liveNamespacesReq)
	if liveNamespacesW.Code != http.StatusOK {
		t.Fatalf("get live namespaces failed: %d body=%s", liveNamespacesW.Code, liveNamespacesW.Body.String())
	}
}

func TestClustersLiveModeWithoutActiveConnection(t *testing.T) {
	r := NewRouter(nil, "live", "")
	req := requestWithRole(http.MethodGet, "/api/v1/clusters", "", "viewer")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d body=%s", http.StatusServiceUnavailable, w.Code, w.Body.String())
	}
}

func TestSwitchCluster(t *testing.T) {
	r := NewRouter(nil, "mock", "")
	req := requestWithRole(http.MethodPost, "/api/v1/clusters/switch", `{"name":"staging-cluster"}`, "admin")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d body=%s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestSwitchClusterRequiresConfirmation(t *testing.T) {
	r := NewRouter(nil, "mock", "")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/clusters/switch", bytes.NewBufferString(`{"name":"staging-cluster"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "admin")
	req.Header.Set("X-User", "tester")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusPreconditionRequired {
		t.Fatalf("expected status %d, got %d body=%s", http.StatusPreconditionRequired, w.Code, w.Body.String())
	}
}

func TestNamespaces(t *testing.T) {
	r := NewRouter(nil, "mock", "")

	listReq, _ := http.NewRequest(http.MethodGet, "/api/v1/namespaces", nil)
	listW := httptest.NewRecorder()
	r.ServeHTTP(listW, listReq)
	if listW.Code != http.StatusOK {
		t.Fatalf("list namespaces failed: %d", listW.Code)
	}

	createReq, _ := http.NewRequest(http.MethodPost, "/api/v1/namespaces", bytes.NewBufferString(`{"name":"qa"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-User-Role", "admin")
	createReq.Header.Set("X-Action-Confirm", "CONFIRM")
	createW := httptest.NewRecorder()
	r.ServeHTTP(createW, createReq)
	if createW.Code != http.StatusCreated {
		t.Fatalf("create namespace failed: %d body=%s", createW.Code, createW.Body.String())
	}

	yamlReq, _ := http.NewRequest(http.MethodGet, "/api/v1/namespaces/qa/yaml", nil)
	yamlW := httptest.NewRecorder()
	r.ServeHTTP(yamlW, yamlReq)
	if yamlW.Code != http.StatusOK {
		t.Fatalf("yaml endpoint failed: %d body=%s", yamlW.Code, yamlW.Body.String())
	}

	delReq, _ := http.NewRequest(http.MethodDelete, "/api/v1/namespaces/qa", nil)
	delReq.Header.Set("X-User-Role", "admin")
	delReq.Header.Set("X-Action-Confirm", "CONFIRM")
	delW := httptest.NewRecorder()
	r.ServeHTTP(delW, delReq)
	if delW.Code != http.StatusNoContent {
		t.Fatalf("delete namespace failed: %d body=%s", delW.Code, delW.Body.String())
	}

	operatorCreateDevReq := requestWithRole(http.MethodPost, "/api/v1/namespaces", `{"name":"dev"}`, "operator")
	operatorCreateDevW := httptest.NewRecorder()
	r.ServeHTTP(operatorCreateDevW, operatorCreateDevReq)
	if operatorCreateDevW.Code != http.StatusConflict {
		t.Fatalf("operator create dev should reach handler and hit conflict, got: %d body=%s", operatorCreateDevW.Code, operatorCreateDevW.Body.String())
	}

	operatorCreateDefaultReq := requestWithRole(http.MethodPost, "/api/v1/namespaces", `{"name":"default"}`, "operator")
	operatorCreateDefaultW := httptest.NewRecorder()
	r.ServeHTTP(operatorCreateDefaultW, operatorCreateDefaultReq)
	if operatorCreateDefaultW.Code != http.StatusForbidden {
		t.Fatalf("operator create default should be forbidden, got: %d body=%s", operatorCreateDefaultW.Code, operatorCreateDefaultW.Body.String())
	}
}

func TestWorkloads(t *testing.T) {
	r := NewRouter(nil, "mock", "")

	deployListReq, _ := http.NewRequest(http.MethodGet, "/api/v1/deployments", nil)
	deployListW := httptest.NewRecorder()
	r.ServeHTTP(deployListW, deployListReq)
	if deployListW.Code != http.StatusOK {
		t.Fatalf("list deployments failed: %d body=%s", deployListW.Code, deployListW.Body.String())
	}

	deployYAMLReq, _ := http.NewRequest(http.MethodGet, "/api/v1/deployments/web-api/yaml", nil)
	deployYAMLW := httptest.NewRecorder()
	r.ServeHTTP(deployYAMLW, deployYAMLReq)
	if deployYAMLW.Code != http.StatusOK {
		t.Fatalf("get deployment yaml failed: %d body=%s", deployYAMLW.Code, deployYAMLW.Body.String())
	}

	updateDeployReq := requestWithRole(http.MethodPut, "/api/v1/deployments/web-api/yaml", `{"yaml":"apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: web-api\n"}`, "operator")
	updateDeployW := httptest.NewRecorder()
	r.ServeHTTP(updateDeployW, updateDeployReq)
	if updateDeployW.Code != http.StatusForbidden {
		t.Fatalf("operator update deployment in default should be forbidden: %d body=%s", updateDeployW.Code, updateDeployW.Body.String())
	}

	updateDeployAdminReq := requestWithRole(http.MethodPut, "/api/v1/deployments/web-api/yaml", `{"yaml":"apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: web-api\n"}`, "admin")
	updateDeployAdminW := httptest.NewRecorder()
	r.ServeHTTP(updateDeployAdminW, updateDeployAdminReq)
	if updateDeployAdminW.Code != http.StatusNoContent {
		t.Fatalf("admin update deployment yaml failed: %d body=%s", updateDeployAdminW.Code, updateDeployAdminW.Body.String())
	}

	updateJobReq := requestWithRole(http.MethodPut, "/api/v1/jobs/db-migrate-20260401/yaml", `{"yaml":"apiVersion: batch/v1\nkind: Job\nmetadata:\n  name: db-migrate-20260401\n"}`, "operator")
	updateJobW := httptest.NewRecorder()
	r.ServeHTTP(updateJobW, updateJobReq)
	if updateJobW.Code != http.StatusForbidden {
		t.Fatalf("operator update default job should be forbidden: %d body=%s", updateJobW.Code, updateJobW.Body.String())
	}

	updatePodDevReq := requestWithRole(http.MethodPut, "/api/v1/pods/task-worker-856ddcf69f-uvwxy/yaml", `{"yaml":"apiVersion: v1\nkind: Pod\nmetadata:\n  name: task-worker-856ddcf69f-uvwxy\n"}`, "operator")
	updatePodDevW := httptest.NewRecorder()
	r.ServeHTTP(updatePodDevW, updatePodDevReq)
	if updatePodDevW.Code != http.StatusNoContent {
		t.Fatalf("operator update pod in dev should pass: %d body=%s", updatePodDevW.Code, updatePodDevW.Body.String())
	}

	podListReq, _ := http.NewRequest(http.MethodGet, "/api/v1/pods", nil)
	podListW := httptest.NewRecorder()
	r.ServeHTTP(podListW, podListReq)
	if podListW.Code != http.StatusOK {
		t.Fatalf("list pods failed: %d body=%s", podListW.Code, podListW.Body.String())
	}

	podLogReq, _ := http.NewRequest(http.MethodGet, "/api/v1/pods/web-api-7bf59f6f9c-abcde/logs", nil)
	podLogW := httptest.NewRecorder()
	r.ServeHTTP(podLogW, podLogReq)
	if podLogW.Code != http.StatusOK {
		t.Fatalf("get pod logs failed: %d body=%s", podLogW.Code, podLogW.Body.String())
	}

	filteredPodLogReq, _ := http.NewRequest(http.MethodGet, "/api/v1/pods/web-api-7bf59f6f9c-abcde/logs?keyword=healthz&matchOnly=true", nil)
	filteredPodLogW := httptest.NewRecorder()
	r.ServeHTTP(filteredPodLogW, filteredPodLogReq)
	if filteredPodLogW.Code != http.StatusOK {
		t.Fatalf("get filtered pod logs failed: %d body=%s", filteredPodLogW.Code, filteredPodLogW.Body.String())
	}

	followPodLogReq, _ := http.NewRequest(http.MethodGet, "/api/v1/pods/web-api-7bf59f6f9c-abcde/logs?follow=true", nil)
	followPodLogW := httptest.NewRecorder()
	r.ServeHTTP(followPodLogW, followPodLogReq)
	if followPodLogW.Code != http.StatusOK {
		t.Fatalf("get follow pod logs failed: %d body=%s", followPodLogW.Code, followPodLogW.Body.String())
	}
	if !bytes.Contains(followPodLogW.Body.Bytes(), []byte("follow refresh tick=")) {
		t.Fatalf("follow pod logs missing refresh marker: %s", followPodLogW.Body.String())
	}

	terminalCapsReq, _ := http.NewRequest(http.MethodGet, "/api/v1/pods/web-api-7bf59f6f9c-abcde/terminal/capabilities", nil)
	terminalCapsW := httptest.NewRecorder()
	r.ServeHTTP(terminalCapsW, terminalCapsReq)
	if terminalCapsW.Code != http.StatusOK {
		t.Fatalf("get terminal capabilities failed: %d body=%s", terminalCapsW.Code, terminalCapsW.Body.String())
	}

	terminalSessionReq := requestWithRole(http.MethodPost, "/api/v1/pods/web-api-7bf59f6f9c-abcde/terminal/sessions", "", "operator")
	terminalSessionW := httptest.NewRecorder()
	r.ServeHTTP(terminalSessionW, terminalSessionReq)
	if terminalSessionW.Code != http.StatusForbidden {
		t.Fatalf("operator terminal session in default should be forbidden: %d body=%s", terminalSessionW.Code, terminalSessionW.Body.String())
	}

	terminalSessionAdminReq := requestWithRole(http.MethodPost, "/api/v1/pods/web-api-7bf59f6f9c-abcde/terminal/sessions", "", "admin")
	terminalSessionAdminW := httptest.NewRecorder()
	r.ServeHTTP(terminalSessionAdminW, terminalSessionAdminReq)
	if terminalSessionAdminW.Code != http.StatusNotImplemented {
		t.Fatalf("admin terminal session placeholder failed: %d body=%s", terminalSessionAdminW.Code, terminalSessionAdminW.Body.String())
	}

	statefulReq, _ := http.NewRequest(http.MethodGet, "/api/v1/statefulsets", nil)
	statefulW := httptest.NewRecorder()
	r.ServeHTTP(statefulW, statefulReq)
	if statefulW.Code != http.StatusOK {
		t.Fatalf("list statefulsets failed: %d body=%s", statefulW.Code, statefulW.Body.String())
	}

	daemonReq, _ := http.NewRequest(http.MethodGet, "/api/v1/daemonsets", nil)
	daemonW := httptest.NewRecorder()
	r.ServeHTTP(daemonW, daemonReq)
	if daemonW.Code != http.StatusOK {
		t.Fatalf("list daemonsets failed: %d body=%s", daemonW.Code, daemonW.Body.String())
	}

	jobReq, _ := http.NewRequest(http.MethodGet, "/api/v1/jobs", nil)
	jobW := httptest.NewRecorder()
	r.ServeHTTP(jobW, jobReq)
	if jobW.Code != http.StatusOK {
		t.Fatalf("list jobs failed: %d body=%s", jobW.Code, jobW.Body.String())
	}

	cronReq, _ := http.NewRequest(http.MethodGet, "/api/v1/cronjobs", nil)
	cronW := httptest.NewRecorder()
	r.ServeHTTP(cronW, cronReq)
	if cronW.Code != http.StatusOK {
		t.Fatalf("list cronjobs failed: %d body=%s", cronW.Code, cronW.Body.String())
	}
}

func TestResourceEndpoints(t *testing.T) {
	r := NewRouter(nil, "mock", "")

	req1, _ := http.NewRequest(http.MethodGet, "/api/v1/services", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("list services failed: %d body=%s", w1.Code, w1.Body.String())
	}

	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/configmaps", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("list configmaps failed: %d body=%s", w2.Code, w2.Body.String())
	}

	req3, _ := http.NewRequest(http.MethodGet, "/api/v1/secrets", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	if w3.Code != http.StatusOK {
		t.Fatalf("list secrets failed: %d body=%s", w3.Code, w3.Body.String())
	}

	req4, _ := http.NewRequest(http.MethodGet, "/api/v1/secrets/web-api-secret", nil)
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, req4)
	if w4.Code != http.StatusOK {
		t.Fatalf("get secret failed: %d body=%s", w4.Code, w4.Body.String())
	}

	req5, _ := http.NewRequest(http.MethodGet, "/api/v1/ingresses", nil)
	w5 := httptest.NewRecorder()
	r.ServeHTTP(w5, req5)
	if w5.Code != http.StatusOK {
		t.Fatalf("list ingresses failed: %d body=%s", w5.Code, w5.Body.String())
	}

	req6, _ := http.NewRequest(http.MethodGet, "/api/v1/ingresses/web-api-ing/services", nil)
	w6 := httptest.NewRecorder()
	r.ServeHTTP(w6, req6)
	if w6.Code != http.StatusOK {
		t.Fatalf("list ingress services failed: %d body=%s", w6.Code, w6.Body.String())
	}

	req7, _ := http.NewRequest(http.MethodGet, "/api/v1/hpas", nil)
	w7 := httptest.NewRecorder()
	r.ServeHTTP(w7, req7)
	if w7.Code != http.StatusOK {
		t.Fatalf("list hpas failed: %d body=%s", w7.Code, w7.Body.String())
	}

	req8, _ := http.NewRequest(http.MethodGet, "/api/v1/hpas/web-api-hpa/target", nil)
	w8 := httptest.NewRecorder()
	r.ServeHTTP(w8, req8)
	if w8.Code != http.StatusOK {
		t.Fatalf("get hpa target failed: %d body=%s", w8.Code, w8.Body.String())
	}

	req9, _ := http.NewRequest(http.MethodGet, "/api/v1/pvs", nil)
	w9 := httptest.NewRecorder()
	r.ServeHTTP(w9, req9)
	if w9.Code != http.StatusOK {
		t.Fatalf("list pvs failed: %d body=%s", w9.Code, w9.Body.String())
	}

	req10, _ := http.NewRequest(http.MethodGet, "/api/v1/pvcs", nil)
	w10 := httptest.NewRecorder()
	r.ServeHTTP(w10, req10)
	if w10.Code != http.StatusOK {
		t.Fatalf("list pvcs failed: %d body=%s", w10.Code, w10.Body.String())
	}

	req11, _ := http.NewRequest(http.MethodGet, "/api/v1/storageclasses", nil)
	w11 := httptest.NewRecorder()
	r.ServeHTTP(w11, req11)
	if w11.Code != http.StatusOK {
		t.Fatalf("list storageclasses failed: %d body=%s", w11.Code, w11.Body.String())
	}

	req12, _ := http.NewRequest(http.MethodGet, "/api/v1/nodes", nil)
	w12 := httptest.NewRecorder()
	r.ServeHTTP(w12, req12)
	if w12.Code != http.StatusOK {
		t.Fatalf("list nodes failed: %d body=%s", w12.Code, w12.Body.String())
	}

	req13, _ := http.NewRequest(http.MethodGet, "/api/v1/nodes/ip-10-10-1-21.ec2.internal", nil)
	w13 := httptest.NewRecorder()
	r.ServeHTTP(w13, req13)
	if w13.Code != http.StatusOK {
		t.Fatalf("get node failed: %d body=%s", w13.Code, w13.Body.String())
	}

	req14, _ := http.NewRequest(http.MethodGet, "/api/v1/nodes/ip-10-10-1-21.ec2.internal/yaml", nil)
	w14 := httptest.NewRecorder()
	r.ServeHTTP(w14, req14)
	if w14.Code != http.StatusOK {
		t.Fatalf("get node yaml failed: %d body=%s", w14.Code, w14.Body.String())
	}

	req15, _ := http.NewRequest(http.MethodGet, "/api/v1/nodes/ip-10-10-1-21.ec2.internal/yaml/download", nil)
	w15 := httptest.NewRecorder()
	r.ServeHTTP(w15, req15)
	if w15.Code != http.StatusOK {
		t.Fatalf("download node yaml failed: %d body=%s", w15.Code, w15.Body.String())
	}
	if cd := w15.Header().Get("Content-Disposition"); cd == "" {
		t.Fatalf("download node yaml missing content-disposition header")
	}
}

func TestRBACAndAudit(t *testing.T) {
	r := NewRouter(nil, "mock", "")

	denyReq := requestWithRole(http.MethodDelete, "/api/v1/namespaces/default", "", "viewer")
	denyW := httptest.NewRecorder()
	r.ServeHTTP(denyW, denyReq)
	if denyW.Code != http.StatusForbidden {
		t.Fatalf("viewer should be forbidden, got %d", denyW.Code)
	}

	auditReq := requestWithRole(http.MethodGet, "/api/v1/audits", "", "admin")
	auditW := httptest.NewRecorder()
	r.ServeHTTP(auditW, auditReq)
	if auditW.Code != http.StatusOK {
		t.Fatalf("admin audit read failed: %d body=%s", auditW.Code, auditW.Body.String())
	}

	filterReq := requestWithRole(http.MethodGet, "/api/v1/audits?method=DELETE&path=/api/v1/namespaces&statusCode=403&limit=1", "", "admin")
	filterW := httptest.NewRecorder()
	r.ServeHTTP(filterW, filterReq)
	if filterW.Code != http.StatusOK {
		t.Fatalf("filtered audit read failed: %d body=%s", filterW.Code, filterW.Body.String())
	}

	clusterManageDenyReq := requestWithRole(http.MethodPost, "/api/v1/clusters/connections/import/token", `{"name":"deny-live","apiServer":"https://deny","bearerToken":"x"}`, "operator")
	clusterManageDenyW := httptest.NewRecorder()
	r.ServeHTTP(clusterManageDenyW, clusterManageDenyReq)
	if clusterManageDenyW.Code != http.StatusForbidden {
		t.Fatalf("operator should be forbidden to import cluster connections, got %d", clusterManageDenyW.Code)
	}
}

func TestAuthUserManagement(t *testing.T) {
	r := NewRouter(nil, "mock", "")

	publicProvidersReq, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/providers/public", nil)
	publicProvidersW := httptest.NewRecorder()
	r.ServeHTTP(publicProvidersW, publicProvidersReq)
	if publicProvidersW.Code != http.StatusOK {
		t.Fatalf("public auth providers should be available: %d body=%s", publicProvidersW.Code, publicProvidersW.Body.String())
	}

	createReq := requestWithRole(http.MethodPost, "/api/v1/auth/users", `{"username":"p601-user","password":"123456","role":"readonly","allowedNamespaces":["dev"]}`, "admin")
	createW := httptest.NewRecorder()
	r.ServeHTTP(createW, createReq)
	if createW.Code != http.StatusServiceUnavailable {
		t.Fatalf("admin create user should return 503 when auth db disabled: %d body=%s", createW.Code, createW.Body.String())
	}

	listReq := requestWithRole(http.MethodGet, "/api/v1/auth/users", "", "admin")
	listW := httptest.NewRecorder()
	r.ServeHTTP(listW, listReq)
	if listW.Code != http.StatusServiceUnavailable {
		t.Fatalf("admin list users should return 503 when auth db disabled: %d body=%s", listW.Code, listW.Body.String())
	}

	listProvidersReq := requestWithRole(http.MethodGet, "/api/v1/auth/providers", "", "admin")
	listProvidersW := httptest.NewRecorder()
	r.ServeHTTP(listProvidersW, listProvidersReq)
	if listProvidersW.Code != http.StatusServiceUnavailable {
		t.Fatalf("admin list providers should return 503 when auth db disabled: %d body=%s", listProvidersW.Code, listProvidersW.Body.String())
	}

	viewerListReq := requestWithRole(http.MethodGet, "/api/v1/auth/users", "", "viewer")
	viewerListW := httptest.NewRecorder()
	r.ServeHTTP(viewerListW, viewerListReq)
	if viewerListW.Code != http.StatusForbidden {
		t.Fatalf("viewer list users should be forbidden: %d body=%s", viewerListW.Code, viewerListW.Body.String())
	}

	disableReq := requestWithRole(http.MethodPatch, "/api/v1/auth/users/p601-user/status", `{"isActive":false}`, "admin")
	disableW := httptest.NewRecorder()
	r.ServeHTTP(disableW, disableReq)
	if disableW.Code != http.StatusServiceUnavailable {
		t.Fatalf("update user status should return 503 when auth db disabled: %d body=%s", disableW.Code, disableW.Body.String())
	}

	updateProfileReq := requestWithRole(http.MethodPatch, "/api/v1/auth/users/p601-user", `{"role":"readonly","allowedNamespaces":["dev"]}`, "admin")
	updateProfileW := httptest.NewRecorder()
	r.ServeHTTP(updateProfileW, updateProfileReq)
	if updateProfileW.Code != http.StatusServiceUnavailable {
		t.Fatalf("update user profile should return 503 when auth db disabled: %d body=%s", updateProfileW.Code, updateProfileW.Body.String())
	}

	resetReq := requestWithRole(http.MethodPost, "/api/v1/auth/users/p601-user/reset-password", `{"password":"654321"}`, "admin")
	resetW := httptest.NewRecorder()
	r.ServeHTTP(resetW, resetReq)
	if resetW.Code != http.StatusServiceUnavailable {
		t.Fatalf("reset user password should return 503 when auth db disabled: %d body=%s", resetW.Code, resetW.Body.String())
	}
}

func TestParseTerminalSessionTTL(t *testing.T) {
	t.Setenv("KM_TERMINAL_SESSION_TTL_SECONDS", "")
	if got := parseTerminalSessionTTL(); got != 120*time.Second {
		t.Fatalf("default ttl mismatch: got=%v want=%v", got, 120*time.Second)
	}

	t.Setenv("KM_TERMINAL_SESSION_TTL_SECONDS", "300")
	if got := parseTerminalSessionTTL(); got != 300*time.Second {
		t.Fatalf("custom ttl mismatch: got=%v want=%v", got, 300*time.Second)
	}

	t.Setenv("KM_TERMINAL_SESSION_TTL_SECONDS", "0")
	if got := parseTerminalSessionTTL(); got != 120*time.Second {
		t.Fatalf("zero ttl should fallback: got=%v", got)
	}

	t.Setenv("KM_TERMINAL_SESSION_TTL_SECONDS", "bad")
	if got := parseTerminalSessionTTL(); got != 120*time.Second {
		t.Fatalf("invalid ttl should fallback: got=%v", got)
	}
}
