package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024 // 512KB
)

// Client represents a connected WebSocket client
type Client struct {
	hub              *Hub
	conn             *websocket.Conn
	send             chan []byte
	UserID           int64
	SessionID        int64
	SessionStudentID int64
	mu               sync.RWMutex
}

// NewClient creates a new WebSocket client
func NewClient(hub *Hub, conn *websocket.Conn, userID, sessionID, sessionStudentID int64) *Client {
	return &Client{
		hub:              hub,
		conn:             conn,
		send:             make(chan []byte, 256),
		UserID:           userID,
		SessionID:        sessionID,
		SessionStudentID: sessionStudentID,
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for user %d: %v", c.UserID, err)
			}
			break
		}

		// Parse the incoming message
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Failed to parse WebSocket message from user %d: %v", c.UserID, err)
			continue
		}

		// Handle the message
		c.hub.handleMessage(c, &msg)
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Drain any queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Send sends a message to the client
func (c *Client) Send(msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	select {
	case c.send <- data:
		return nil
	default:
		return ErrClientBufferFull
	}
}

// SendTimerSync sends a timer synchronization message to the client
func (c *Client) SendTimerSync(remainingSeconds int64, endsAt int64) error {
	msg, err := NewMessage(TypeTimerSync, TimerSync{
		RemainingSeconds: remainingSeconds,
		EndsAt:           endsAt,
	})
	if err != nil {
		return err
	}
	return c.Send(msg)
}

// SendSessionEnd sends a session end message to the client
func (c *Client) SendSessionEnd(reason string) error {
	msg, err := NewMessage(TypeSessionEnd, SessionEnd{Reason: reason})
	if err != nil {
		return err
	}
	return c.Send(msg)
}

// SendError sends an error message to the client
func (c *Client) SendError(code, message string) error {
	msg, err := NewMessage(TypeError, ErrorResponse{Code: code, Message: message})
	if err != nil {
		return err
	}
	return c.Send(msg)
}

// Close closes the client connection
func (c *Client) Close() {
	close(c.send)
}

// ErrClientBufferFull is returned when the client's send buffer is full
var ErrClientBufferFull = &ClientError{Message: "client send buffer full"}

// ClientError represents a client-related error
type ClientError struct {
	Message string
}

func (e *ClientError) Error() string {
	return e.Message
}
