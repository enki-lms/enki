package exams

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/enki/daemon/internal/auth"
	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/enki/daemon/internal/websocket"
	"github.com/gin-gonic/gin"
	gorillaws "github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now; tighten in production
	},
}

// Handler handles exam-related API routes
type Handler struct {
	queries    *sqlc.Queries
	middleware *auth.Middleware
	hub        *websocket.Hub
}

// NewHandler creates a new exam handler
func NewHandler(queries *sqlc.Queries, middleware *auth.Middleware, hub *websocket.Hub) *Handler {
	return &Handler{
		queries:    queries,
		middleware: middleware,
		hub:        hub,
	}
}

// CreateExamSessionRequest represents the request body for creating an exam session
type CreateExamSessionRequest struct {
	ProblemGroupType string  `json:"problem_group_type" binding:"required,oneof=comp_sci quiz"`
	ProblemGroupID   int64   `json:"problem_group_id" binding:"required"`
	DurationMinutes  int32   `json:"duration_minutes" binding:"required,min=1"`
	StudentIDs       []int64 `json:"student_ids" binding:"required,min=1"`
}

// CreateExamSession creates a new exam session
// POST /api/exams/sessions
func (h *Handler) CreateExamSession(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateExamSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the session
	session, err := h.queries.CreateExamSession(c.Request.Context(), sqlc.CreateExamSessionParams{
		ProblemGroupType: req.ProblemGroupType,
		ProblemGroupID:   req.ProblemGroupID,
		OpenedBy:         claims.UserID,
		DurationMinutes:  req.DurationMinutes,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create exam session"})
		return
	}

	// Add students to the session
	for _, studentID := range req.StudentIDs {
		_, err := h.queries.CreateExamSessionStudent(c.Request.Context(), sqlc.CreateExamSessionStudentParams{
			SessionID: session.ID,
			UserID:    studentID,
		})
		if err != nil {
			// Log but continue - some students might already be added
			continue
		}
	}

	c.JSON(http.StatusCreated, session)
}

// GetExamSession returns an exam session by ID
// GET /api/exams/sessions/:id
func (h *Handler) GetExamSession(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	session, err := h.queries.GetExamSession(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get session"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// ListExamSessionStudents returns all students in an exam session
// GET /api/exams/sessions/:id/students
func (h *Handler) ListExamSessionStudents(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	students, err := h.queries.ListExamSessionStudents(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list students"})
		return
	}

	c.JSON(http.StatusOK, students)
}

// StartExamSession starts an exam session
// POST /api/exams/sessions/:id/start
func (h *Handler) StartExamSession(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	session, err := h.queries.StartExamSession(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start session"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// EndExamSession ends an exam session
// POST /api/exams/sessions/:id/end
func (h *Handler) EndExamSession(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	session, err := h.queries.EndExamSession(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to end session"})
		return
	}

	// Notify all connected students
	h.hub.EndSession(id, "teacher_ended")

	c.JSON(http.StatusOK, session)
}

// DiscontinueStudent discontinues a student from an exam session
// POST /api/exams/sessions/:id/students/:studentId/discontinue
func (h *Handler) DiscontinueStudent(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	studentID, err := strconv.ParseInt(c.Param("studentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	// Find the session student record
	sessionStudent, err := h.queries.GetExamSessionStudentBySessionAndUser(c.Request.Context(), sqlc.GetExamSessionStudentBySessionAndUserParams{
		SessionID: sessionID,
		UserID:    studentID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not in session"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find student"})
		return
	}

	// Discontinue the student
	_, err = h.queries.DiscontinueExamSessionStudent(c.Request.Context(), sessionStudent.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to discontinue student"})
		return
	}

	// Disconnect them from WebSocket
	h.hub.EndSessionForStudent(sessionID, studentID, "discontinued")

	c.JSON(http.StatusOK, gin.H{"message": "student discontinued"})
}

// ListActiveExamSessionsForUser returns active exam sessions for the current user
// GET /api/exams/sessions/active
func (h *Handler) ListActiveExamSessionsForUser(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sessions, err := h.queries.ListActiveExamSessionsForUser(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list sessions"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// ListExamSessionsByTeacher returns exam sessions created by the current teacher
// GET /api/exams/sessions
func (h *Handler) ListExamSessionsByTeacher(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sessions, err := h.queries.ListExamSessionsByTeacher(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list sessions"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// HandleExamWebSocket handles WebSocket connections for exam sessions
// GET /ws/exam/:sessionId
func (h *Handler) HandleExamWebSocket(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sessionID, err := strconv.ParseInt(c.Param("sessionId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	// Get the session
	session, err := h.queries.GetExamSession(c.Request.Context(), sessionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get session"})
		return
	}

	// Check session is active
	if session.Status != sqlc.ExamSessionStatusActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session not active"})
		return
	}

	// Get the student's session record
	sessionStudent, err := h.queries.GetExamSessionStudentBySessionAndUser(c.Request.Context(), sqlc.GetExamSessionStudentBySessionAndUserParams{
		SessionID: sessionID,
		UserID:    claims.UserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not assigned to this exam"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check assignment"})
		return
	}

	// Check student status
	if sessionStudent.Status == sqlc.ExamStudentStatusSubmitted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already submitted"})
		return
	}
	if sessionStudent.Status == sqlc.ExamStudentStatusDiscontinued {
		c.JSON(http.StatusForbidden, gin.H{"error": "discontinued from exam"})
		return
	}

	// Calculate end time - students get full duration from when they join
	var endsAt time.Time
	if sessionStudent.JoinedAt.Valid {
		// Already joined before - resume with same end time
		endsAt = sessionStudent.EndsAt.Time
	} else {
		// First time joining - set end time
		endsAt = time.Now().Add(time.Duration(session.DurationMinutes) * time.Minute)
		sessionStudent, err = h.queries.JoinExamSession(c.Request.Context(), sqlc.JoinExamSessionParams{
			ID:     sessionStudent.ID,
			EndsAt: pgtype.Timestamp{Time: endsAt, Valid: true},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to join session"})
			return
		}
	}

	// Check if time already expired
	if time.Now().After(endsAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exam time expired"})
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// Create client
	client := websocket.NewClient(h.hub, conn, claims.UserID, sessionID, sessionStudent.ID)

	// Register client
	h.hub.Register(client)

	// Send initial timer sync
	remainingSeconds := int64(time.Until(endsAt).Seconds())
	client.SendTimerSync(remainingSeconds, endsAt.Unix())

	// Also send joined response with exam info
	joinedMsg, _ := websocket.NewMessage(websocket.TypeJoined, websocket.JoinedResponse{
		SessionID:        sessionID,
		RemainingSeconds: remainingSeconds,
		EndsAt:           endsAt.Unix(),
		ProblemGroupType: session.ProblemGroupType,
		ProblemGroupID:   session.ProblemGroupID,
	})
	client.Send(joinedMsg)

	// Start read/write pumps
	go client.WritePump()
	go client.ReadPump()
}

// RegisterRoutes registers exam routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware *auth.Middleware) {
	exams := rg.Group("/exams")
	{
		// Student endpoints
		exams.GET("/sessions/active", authMiddleware.AuthRequired(), h.ListActiveExamSessionsForUser)

		// Teacher endpoints
		exams.GET("/sessions", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.ListExamSessionsByTeacher)
		exams.POST("/sessions", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.CreateExamSession)
		exams.GET("/sessions/:id", authMiddleware.AuthRequired(), h.GetExamSession)
		exams.GET("/sessions/:id/students", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.ListExamSessionStudents)
		exams.POST("/sessions/:id/start", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.StartExamSession)
		exams.POST("/sessions/:id/end", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.EndExamSession)
		exams.POST("/sessions/:id/students/:studentId/discontinue", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.DiscontinueStudent)
	}
}

// RegisterWebSocketRoutes registers WebSocket routes
func (h *Handler) RegisterWebSocketRoutes(r *gin.Engine, authMiddleware *auth.Middleware) {
	r.GET("/ws/exam/:sessionId", authMiddleware.AuthRequired(), h.HandleExamWebSocket)
}
