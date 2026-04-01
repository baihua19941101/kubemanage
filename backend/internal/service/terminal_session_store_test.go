package service

import (
	"testing"
	"time"
)

func TestTerminalSessionStoreConsumeOnce(t *testing.T) {
	store := NewTerminalSessionStore(2 * time.Minute)
	session := store.Create("pod-a", "container-a", "alice", "admin")

	got, err := store.Consume(session.ID, "pod-a")
	if err != nil {
		t.Fatalf("first consume should succeed: %v", err)
	}
	if got.ID != session.ID {
		t.Fatalf("unexpected session id: got=%s want=%s", got.ID, session.ID)
	}
	if got.PodName != "pod-a" {
		t.Fatalf("unexpected pod name: %s", got.PodName)
	}
	if got.Container != "container-a" {
		t.Fatalf("unexpected container: %s", got.Container)
	}

	_, err = store.Consume(session.ID, "pod-a")
	if err != ErrTerminalSessionNotFound {
		t.Fatalf("second consume should fail with not found, got: %v", err)
	}
}

func TestTerminalSessionStoreExpire(t *testing.T) {
	store := NewTerminalSessionStore(20 * time.Millisecond)
	session := store.Create("pod-a", "container-a", "alice", "admin")
	time.Sleep(40 * time.Millisecond)

	_, err := store.Consume(session.ID, "pod-a")
	if err != ErrTerminalSessionExpired {
		t.Fatalf("expired session should return expired, got: %v", err)
	}
}

func TestTerminalSessionStorePodMismatch(t *testing.T) {
	store := NewTerminalSessionStore(2 * time.Minute)
	session := store.Create("pod-a", "container-a", "alice", "admin")

	_, err := store.Consume(session.ID, "pod-b")
	if err != ErrTerminalSessionNotFound {
		t.Fatalf("pod mismatch should fail with not found, got: %v", err)
	}
}

func TestTerminalSessionStoreTTL(t *testing.T) {
	store := NewTerminalSessionStore(90 * time.Second)
	if got := store.TTL(); got != 90*time.Second {
		t.Fatalf("ttl mismatch: got=%v", got)
	}

	fallbackStore := NewTerminalSessionStore(0)
	if got := fallbackStore.TTL(); got != 2*time.Minute {
		t.Fatalf("fallback ttl mismatch: got=%v", got)
	}
}
