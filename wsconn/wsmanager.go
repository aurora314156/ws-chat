package wsconn

import (
	"context"
	"sync"
	"ws-chat/logger"

	"github.com/gorilla/websocket"
)

type WSManager struct {
	conns map[*websocket.Conn]struct{}
	mu    sync.Mutex
}

func New() *WSManager {
	return &WSManager{
		conns: make(map[*websocket.Conn]struct{}),
	}
}

// add new ws connection
func (m *WSManager) Add(conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.conns[conn] = struct{}{}
}

// remove ws connection
func (m *WSManager) Remove(conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.conns, conn)
}

// gracefully close all ws connections
func (m *WSManager) CloseAll(ctx context.Context) {
	m.mu.Lock()
	conns := make([]*websocket.Conn, 0, len(m.conns))
	for c := range m.conns {
		conns = append(conns, c)
	}
	m.mu.Unlock()

	done := make(chan struct{})
	go func() {
		for _, c := range conns {
			_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "live chat server shutting down"))
			c.Close()
		}
		close(done)
	}()

	select {
	case <-ctx.Done():
		logger.Info("[WSManager] Timeout closing WebSocket connections")
	case <-done:
		logger.Info("[WSManager] All WebSocket connections closed")
	}
}

// Broadcast sends a message to all active WebSocket connections
func (m *WSManager) Broadcast(msg map[string]any) {
	m.mu.Lock()
	conns := make([]*websocket.Conn, 0, len(m.conns))
	for c := range m.conns {
		conns = append(conns, c)
	}
	m.mu.Unlock()

	for _, c := range conns {
		if err := c.WriteJSON(msg); err != nil {
			logger.Error("[WSManager] Broadcast error:", err)
			m.Remove(c)
			c.Close()
		}
	}
}
