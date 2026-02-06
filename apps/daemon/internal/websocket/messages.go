package websocket

import "encoding/json"

// MessageType constants for WebSocket messages
const (
	TypeTimerSync    = "timer_sync"
	TypeSaveWork     = "save_work"
	TypeWorkSaved    = "work_saved"
	TypeSubmit       = "submit"
	TypeSessionEnd   = "session_end"
	TypeError        = "error"
	TypeJoined       = "joined"
	TypeDisconnected = "disconnected"
)

// Message is the base WebSocket message structure
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// TimerSync is sent to clients to synchronize the exam timer
type TimerSync struct {
	RemainingSeconds int64 `json:"remaining_seconds"`
	EndsAt           int64 `json:"ends_at"` // Unix timestamp
}

// SaveWork is sent by clients to save their current work
type SaveWork struct {
	ProblemID   int64  `json:"problem_id"`
	ProblemType string `json:"problem_type"` // "comp_sci" or "quiz"
	Answer      string `json:"answer"`       // Code or JSON answer
}

// WorkSaved confirms that work was saved successfully
type WorkSaved struct {
	ProblemID   int64  `json:"problem_id"`
	ProblemType string `json:"problem_type"`
	SavedAt     int64  `json:"saved_at"` // Unix timestamp
}

// Submit is sent by clients to submit their exam
type Submit struct {
	// Empty - just the message type is enough
}

// SessionEnd is broadcast when the exam session ends
type SessionEnd struct {
	Reason string `json:"reason"` // "time_expired", "teacher_ended", "submitted"
}

// JoinedResponse is sent when a student successfully joins an exam
type JoinedResponse struct {
	SessionID        int64  `json:"session_id"`
	RemainingSeconds int64  `json:"remaining_seconds"`
	EndsAt           int64  `json:"ends_at"`
	ProblemGroupType string `json:"problem_group_type"`
	ProblemGroupID   int64  `json:"problem_group_id"`
}

// ErrorResponse is sent when an error occurs
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewMessage creates a new message with the given type and payload
func NewMessage(msgType string, payload interface{}) (*Message, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &Message{
		Type:    msgType,
		Payload: data,
	}, nil
}

// ParsePayload parses the message payload into the given struct
func (m *Message) ParsePayload(v interface{}) error {
	return json.Unmarshal(m.Payload, v)
}
