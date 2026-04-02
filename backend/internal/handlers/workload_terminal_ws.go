package handlers

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"

	"kubeManage/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var terminalUpgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

type wsTextWriter struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (w *wsTextWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.conn.WriteMessage(websocket.TextMessage, p); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (h *WorkloadHandler) TerminalWebSocket(c *gin.Context) {
	if h.adapterMode == "mock" || h.liveWorkloadSvc == nil || h.terminalSessions == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "terminal websocket not enabled in mock mode"})
		return
	}
	podName := c.Param("name")
	sessionID := c.Query("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sessionId is required"})
		return
	}
	session, err := h.terminalSessions.Get(sessionID, podName)
	if err != nil {
		switch err {
		case service.ErrTerminalSessionNotFound:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "terminal session invalid"})
		case service.ErrTerminalSessionExpired:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "terminal session expired"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	// Bind websocket attach to the same caller that created the session.
	requestUser := strings.TrimSpace(c.GetString("km_user"))
	requestRole := strings.TrimSpace(c.GetString("km_role"))
	if requestUser == "" || requestUser == "demo-user" {
		queryUser := strings.TrimSpace(c.Query("user"))
		if queryUser != "" {
			requestUser = queryUser
		}
	}
	if requestRole == "" || requestRole == "readonly" {
		queryRole := strings.TrimSpace(c.Query("role"))
		if queryRole != "" {
			requestRole = queryRole
		}
	}
	if requestUser != "" && session.User != "" && requestUser != session.User {
		c.JSON(http.StatusForbidden, gin.H{"error": "terminal session owner mismatch"})
		return
	}
	if requestRole != "" && session.Role != "" && requestRole != session.Role {
		c.JSON(http.StatusForbidden, gin.H{"error": "terminal session owner mismatch"})
		return
	}
	session, err = h.terminalSessions.Consume(sessionID, podName)
	if err != nil {
		switch err {
		case service.ErrTerminalSessionNotFound:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "terminal session invalid"})
		case service.ErrTerminalSessionExpired:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "terminal session expired"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	conn, err := terminalUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	command := c.QueryArray("command")
	if len(command) == 0 {
		command = []string{"sh"}
	}

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	stdinReader, stdinWriter := ioPipe()
	defer stdinWriter.Close()

	writer := &wsTextWriter{conn: conn}
	execDone := make(chan error, 1)
	go func() {
		execDone <- h.liveWorkloadSvc.ExecuteTerminal(ctx, session.PodName, session.Container, command, stdinReader, writer, writer, true)
	}()

	readDone := make(chan struct{})
	go func() {
		defer close(readDone)
		for {
			msgType, data, readErr := conn.ReadMessage()
			if readErr != nil {
				return
			}
			if msgType == websocket.TextMessage || msgType == websocket.BinaryMessage {
				if _, writeErr := stdinWriter.Write(data); writeErr != nil {
					return
				}
			}
		}
	}()

	select {
	case err = <-execDone:
		if err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("terminal error: "+err.Error()))
		}
	case <-readDone:
		cancel()
	case <-ctx.Done():
	}
}

func ioPipe() (*io.PipeReader, *io.PipeWriter) {
	return io.Pipe()
}
