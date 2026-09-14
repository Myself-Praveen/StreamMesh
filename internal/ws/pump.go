package ws

import (
	"bytes"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// ReadPump pumps messages from the websocket connection to the hub.
func (c *Connection) ReadPump(unregister func(*Connection), processMessage func(*Connection, []byte)) {
	defer func() {
		unregister(c)
		c.Conn.Close()
	}()

	hbConfig := GetHeartbeatConfig()

	c.Conn.SetReadLimit(hbConfig.MaxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(hbConfig.PongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(hbConfig.PongWait))
		c.UpdateLastPing()
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Log.Error("websocket unexpected close error", zap.Error(err), zap.String("id", c.ID))
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		processMessage(c, message)
	}
}

// WritePump pumps messages from the hub to the websocket connection.
func (c *Connection) WritePump() {
	hbConfig := GetHeartbeatConfig()
	ticker := time.NewTicker(hbConfig.PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(hbConfig.WriteWait))
			if !ok {
				// The hub closed the channel.
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(hbConfig.WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
