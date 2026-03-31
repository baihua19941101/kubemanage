package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	req, err := http.NewRequest(http.MethodPost, "/api/v1/clusters/switch", bytes.NewBufferString(`{"name":"staging-cluster"}`))
	if err != nil {
		t.Fatalf("build request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

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
	delW := httptest.NewRecorder()
	r.ServeHTTP(delW, delReq)
	if delW.Code != http.StatusNoContent {
		t.Fatalf("delete namespace failed: %d body=%s", delW.Code, delW.Body.String())
	}
}
