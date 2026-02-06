package websocket

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// MessageHandler is called when a client sends a message
type MessageHandler func(client *Client, msg *Message)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients by session ID
	sessions map[int64]map[int64]*Client // sessionID -> userID -> client

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Database queries for saving work
	queries *sqlc.Queries

	// Mutex for sessions map
	mu sync.RWMutex

	// Custom message handlers
	handlers map[string]MessageHandler
}

// NewHub creates a new Hub
func NewHub(queries *sqlc.Queries) *Hub {
	h := &Hub{
		sessions:   make(map[int64]map[int64]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		queries:    queries,
		handlers:   make(map[string]MessageHandler),
	}

	// Register default handlers
	h.handlers[TypeSaveWork] = h.handleSaveWork
	h.handlers[TypeSubmit] = h.handleSubmit

	return h
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.sessions[client.SessionID] == nil {
				h.sessions[client.SessionID] = make(map[int64]*Client)
			}
			h.sessions[client.SessionID][client.UserID] = client
			h.mu.Unlock()
			log.Printf("Client registered: user=%d, session=%d", client.UserID, client.SessionID)

		case client := <-h.unregister:
			h.mu.Lock()
			if users, ok := h.sessions[client.SessionID]; ok {
				if _, ok := users[client.UserID]; ok {
					delete(users, client.UserID)
					client.Close()
					if len(users) == 0 {
						delete(h.sessions, client.SessionID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client unregistered: user=%d, session=%d", client.UserID, client.SessionID)
		}
	}
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// GetClient returns a client by session and user ID
func (h *Hub) GetClient(sessionID, userID int64) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if users, ok := h.sessions[sessionID]; ok {
		return users[userID]
	}
	return nil
}

// GetSessionClients returns all clients in a session
func (h *Hub) GetSessionClients(sessionID int64) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients := make([]*Client, 0)
	if users, ok := h.sessions[sessionID]; ok {
		for _, client := range users {
			clients = append(clients, client)
		}
	}
	return clients
}

// BroadcastToSession sends a message to all clients in a session
func (h *Hub) BroadcastToSession(sessionID int64, msg *Message) {
	clients := h.GetSessionClients(sessionID)
	for _, client := range clients {
		if err := client.Send(msg); err != nil {
			log.Printf("Failed to send message to user %d: %v", client.UserID, err)
		}
	}
}

// EndSessionForStudent sends a session end message and disconnects a student
func (h *Hub) EndSessionForStudent(sessionID, userID int64, reason string) {
	client := h.GetClient(sessionID, userID)
	if client != nil {
		client.SendSessionEnd(reason)
		h.Unregister(client)
	}
}

// EndSession ends a session for all students
func (h *Hub) EndSession(sessionID int64, reason string) {
	clients := h.GetSessionClients(sessionID)
	for _, client := range clients {
		client.SendSessionEnd(reason)
	}
	// Let clients disconnect gracefully, then unregister
	time.AfterFunc(time.Second, func() {
		for _, client := range clients {
			h.Unregister(client)
		}
	})
}

// handleMessage routes incoming messages to the appropriate handler
func (h *Hub) handleMessage(client *Client, msg *Message) {
	if handler, ok := h.handlers[msg.Type]; ok {
		handler(client, msg)
	} else {
		log.Printf("Unknown message type from user %d: %s", client.UserID, msg.Type)
		client.SendError("unknown_message_type", "Unknown message type: "+msg.Type)
	}
}

// handleSaveWork saves a student's work in progress
func (h *Hub) handleSaveWork(client *Client, msg *Message) {
	var payload SaveWork
	if err := msg.ParsePayload(&payload); err != nil {
		client.SendError("invalid_payload", "Invalid save_work payload")
		return
	}

	ctx := context.Background()
	_, err := h.queries.UpsertExamWorkInProgress(ctx, sqlc.UpsertExamWorkInProgressParams{
		SessionStudentID: client.SessionStudentID,
		ProblemID:        payload.ProblemID,
		ProblemType:      payload.ProblemType,
		CurrentAnswer:    payload.Answer,
	})

	if err != nil {
		log.Printf("Failed to save work for user %d: %v", client.UserID, err)
		client.SendError("save_failed", "Failed to save work")
		return
	}

	// Confirm save
	response, _ := NewMessage(TypeWorkSaved, WorkSaved{
		ProblemID:   payload.ProblemID,
		ProblemType: payload.ProblemType,
		SavedAt:     time.Now().Unix(),
	})
	client.Send(response)
}

// handleSubmit handles a student submitting their exam
func (h *Hub) handleSubmit(client *Client, msg *Message) {
	ctx := context.Background()

	// Mark as submitted
	_, err := h.queries.SubmitExamSession(ctx, sqlc.SubmitExamSessionParams{
		ID:            client.SessionStudentID,
		AutoSubmitted: pgtype.Bool{Bool: false, Valid: true},
	})

	if err != nil {
		log.Printf("Failed to submit exam for user %d: %v", client.UserID, err)
		client.SendError("submit_failed", "Failed to submit exam")
		return
	}

	// Send session end message
	client.SendSessionEnd("submitted")
	h.Unregister(client)
}

// SendTimerSyncToClient sends a timer sync to a specific client
func (h *Hub) SendTimerSyncToClient(sessionID, userID int64, remainingSeconds, endsAt int64) {
	client := h.GetClient(sessionID, userID)
	if client != nil {
		client.SendTimerSync(remainingSeconds, endsAt)
	}
}
