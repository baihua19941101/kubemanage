package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func setupTerminalWSRouter(handler *WorkloadHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		user := c.GetHeader("X-User")
		if user == "" {
			user = "demo-user"
		}
		role := c.GetHeader("X-User-Role")
		if role == "" {
			role = "viewer"
		}
		c.Set("km_user", user)
		c.Set("km_role", role)
		c.Next()
	})
	r.GET("/api/v1/pods/:name/terminal/ws", handler.TerminalWebSocket)
	return r
}

func TestTerminalWebSocketMockModeDisabled(t *testing.T) {
	handler := NewWorkloadHandler(service.NewWorkloadService(), service.NewLiveWorkloadReader(nil), service.NewTerminalSessionStore(2*time.Minute), "mock")
	r := setupTerminalWSRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pods/pod-a/terminal/ws?sessionId=any", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Fatalf("unexpected status: %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "not enabled") {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestTerminalWebSocketSessionRequired(t *testing.T) {
	handler := NewWorkloadHandler(service.NewWorkloadService(), service.NewLiveWorkloadReader(nil), service.NewTerminalSessionStore(2*time.Minute), "live")
	r := setupTerminalWSRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pods/pod-a/terminal/ws", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "sessionId is required") {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestTerminalWebSocketInvalidSession(t *testing.T) {
	handler := NewWorkloadHandler(service.NewWorkloadService(), service.NewLiveWorkloadReader(nil), service.NewTerminalSessionStore(2*time.Minute), "live")
	r := setupTerminalWSRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pods/pod-a/terminal/ws?sessionId=invalid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "terminal session invalid") {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestTerminalWebSocketPodMismatch(t *testing.T) {
	store := service.NewTerminalSessionStore(2 * time.Minute)
	session := store.Create("pod-a", "container-a", "alice", "admin")
	handler := NewWorkloadHandler(service.NewWorkloadService(), service.NewLiveWorkloadReader(nil), store, "live")
	r := setupTerminalWSRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pods/pod-b/terminal/ws?sessionId="+session.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "terminal session invalid") {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestTerminalWebSocketExpiredSession(t *testing.T) {
	store := service.NewTerminalSessionStore(20 * time.Millisecond)
	session := store.Create("pod-a", "container-a", "alice", "admin")
	time.Sleep(40 * time.Millisecond)
	handler := NewWorkloadHandler(service.NewWorkloadService(), service.NewLiveWorkloadReader(nil), store, "live")
	r := setupTerminalWSRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pods/pod-a/terminal/ws?sessionId="+session.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "terminal session expired") {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestTerminalWebSocketOwnerMismatch(t *testing.T) {
	store := service.NewTerminalSessionStore(2 * time.Minute)
	session := store.Create("pod-a", "container-a", "alice", "admin")
	handler := NewWorkloadHandler(service.NewWorkloadService(), service.NewLiveWorkloadReader(nil), store, "live")
	r := setupTerminalWSRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pods/pod-a/terminal/ws?sessionId="+session.ID, nil)
	req.Header.Set("X-User", "bob")
	req.Header.Set("X-User-Role", "admin")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("unexpected status: %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "owner mismatch") {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}
