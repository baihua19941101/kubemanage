package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
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
	return req
}

func TestHealthz(t *testing.T) {
	r := NewRouter(nil)
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
	r := NewRouter(nil)
	req, err := http.NewRequest(http.MethodGet, "/api/v1/clusters", nil)
	if err != nil {
		t.Fatalf("build request failed: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSwitchCluster(t *testing.T) {
	r := NewRouter(nil)
	req := requestWithRole(http.MethodPost, "/api/v1/clusters/switch", `{"name":"staging-cluster"}`, "admin")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d body=%s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestNamespaces(t *testing.T) {
	r := NewRouter(nil)

	listReq, _ := http.NewRequest(http.MethodGet, "/api/v1/namespaces", nil)
	listW := httptest.NewRecorder()
	r.ServeHTTP(listW, listReq)
	if listW.Code != http.StatusOK {
		t.Fatalf("list namespaces failed: %d", listW.Code)
	}

	createReq, _ := http.NewRequest(http.MethodPost, "/api/v1/namespaces", bytes.NewBufferString(`{"name":"qa"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-User-Role", "admin")
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
	delW := httptest.NewRecorder()
	r.ServeHTTP(delW, delReq)
	if delW.Code != http.StatusNoContent {
		t.Fatalf("delete namespace failed: %d body=%s", delW.Code, delW.Body.String())
	}
}

func TestWorkloads(t *testing.T) {
	r := NewRouter(nil)

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
	if updateDeployW.Code != http.StatusNoContent {
		t.Fatalf("update deployment yaml failed: %d body=%s", updateDeployW.Code, updateDeployW.Body.String())
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
	r := NewRouter(nil)

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
}

func TestRBACAndAudit(t *testing.T) {
	r := NewRouter(nil)

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
}
