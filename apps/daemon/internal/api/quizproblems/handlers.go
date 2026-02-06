package quizproblems

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/enki/daemon/internal/auth"
	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Handler handles quiz problem-related API routes
type Handler struct {
	queries    *sqlc.Queries
	middleware *auth.Middleware
}

// NewHandler creates a new quiz problem handler
func NewHandler(queries *sqlc.Queries, middleware *auth.Middleware) *Handler {
	return &Handler{
		queries:    queries,
		middleware: middleware,
	}
}

// getGroupInstitution gets the institution for a quiz problem group via its course
func (h *Handler) getGroupInstitution(c *gin.Context, groupID int64) (string, error) {
	group, err := h.queries.GetQuizProblemGroup(c.Request.Context(), groupID)
	if err != nil {
		return "", err
	}
	course, err := h.queries.GetCourse(c.Request.Context(), group.CourseID)
	if err != nil {
		return "", err
	}
	return course.Institution, nil
}

// getProblemInstitution gets the institution for a quiz problem
func (h *Handler) getProblemInstitution(c *gin.Context, problemID int64) (string, error) {
	problem, err := h.queries.GetQuizProblem(c.Request.Context(), problemID)
	if err != nil {
		return "", err
	}
	return h.getGroupInstitution(c, problem.GroupID)
}

// ListQuizProblems returns quiz problems for a group
// GET /api/quiz-groups/:id/problems
func (h *Handler) ListQuizProblems(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	// Verify group belongs to user's institution
	institution, err := h.getGroupInstitution(c, groupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "quiz group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get group"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	problems, err := h.queries.ListQuizProblemsByGroup(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list quiz problems"})
		return
	}

	c.JSON(http.StatusOK, problems)
}

// QuizProblemResponse includes problem and its options
type QuizProblemResponse struct {
	Problem sqlc.QuizProblem         `json:"problem"`
	Options []sqlc.QuizProblemOption `json:"options,omitempty"`
}

// GetQuizProblem returns a quiz problem by ID with its options
// GET /api/quiz-problems/:id
func (h *Handler) GetQuizProblem(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid problem id"})
		return
	}

	problem, err := h.queries.GetQuizProblem(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "quiz problem not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get quiz problem"})
		return
	}

	// Check institution access
	institution, err := h.getGroupInstitution(c, problem.GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Get options for MCQ/true-false problems
	var options []sqlc.QuizProblemOption
	if problem.ProblemType == sqlc.QuizProblemTypeTrueFalse ||
		problem.ProblemType == sqlc.QuizProblemTypeMcqSingle ||
		problem.ProblemType == sqlc.QuizProblemTypeMcqMulti {
		options, err = h.queries.ListQuizProblemOptionsByProblem(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get options"})
			return
		}

		// For students, hide which options are correct
		if !h.middleware.IsTeacherOrAdmin(claims) {
			for i := range options {
				options[i].IsCorrect = false
			}
		}
	}

	c.JSON(http.StatusOK, QuizProblemResponse{
		Problem: problem,
		Options: options,
	})
}

// OptionRequest represents an option in create/update requests
type OptionRequest struct {
	OptionText   string `json:"option_text" binding:"required"`
	IsCorrect    bool   `json:"is_correct"`
	DisplayOrder int32  `json:"display_order"`
}

// CreateQuizProblemRequest represents the request body for creating a quiz problem
type CreateQuizProblemRequest struct {
	ProblemType   string          `json:"problem_type" binding:"required,oneof=open_ended true_false mcq_single mcq_multi fill_blank"`
	Name          string          `json:"name" binding:"required"`
	Description   string          `json:"description" binding:"required"`
	ProblemText   string          `json:"problem_text" binding:"required"`
	Points        int32           `json:"points"`
	CorrectAnswer *string         `json:"correct_answer,omitempty"` // For fill_blank
	Options       []OptionRequest `json:"options,omitempty"`        // For MCQ/true-false
}

// CreateQuizProblem creates a new quiz problem
// POST /api/quiz-groups/:id/problems
func (h *Handler) CreateQuizProblem(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	var req CreateQuizProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify group belongs to user's institution
	institution, err := h.getGroupInstitution(c, groupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "quiz group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get group"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Set default points
	points := req.Points
	if points == 0 {
		points = 1
	}

	// Build correct answer field
	var correctAnswer pgtype.Text
	if req.CorrectAnswer != nil {
		correctAnswer = pgtype.Text{String: *req.CorrectAnswer, Valid: true}
	}

	problem, err := h.queries.CreateQuizProblem(c.Request.Context(), sqlc.CreateQuizProblemParams{
		GroupID:       groupID,
		ProblemType:   sqlc.QuizProblemType(req.ProblemType),
		Name:          req.Name,
		Description:   req.Description,
		ProblemText:   req.ProblemText,
		Points:        points,
		CorrectAnswer: correctAnswer,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create quiz problem"})
		return
	}

	// Create options if provided
	var options []sqlc.QuizProblemOption
	for _, opt := range req.Options {
		option, err := h.queries.CreateQuizProblemOption(c.Request.Context(), sqlc.CreateQuizProblemOptionParams{
			ProblemID:    problem.ID,
			OptionText:   opt.OptionText,
			IsCorrect:    opt.IsCorrect,
			DisplayOrder: opt.DisplayOrder,
		})
		if err != nil {
			// Clean up: delete problem if options fail
			_ = h.queries.DeleteQuizProblem(c.Request.Context(), problem.ID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create option"})
			return
		}
		options = append(options, option)
	}

	c.JSON(http.StatusCreated, QuizProblemResponse{
		Problem: problem,
		Options: options,
	})
}

// UpdateQuizProblemRequest represents the request body for updating a quiz problem
type UpdateQuizProblemRequest struct {
	ProblemType   string          `json:"problem_type" binding:"required,oneof=open_ended true_false mcq_single mcq_multi fill_blank"`
	Name          string          `json:"name" binding:"required"`
	Description   string          `json:"description" binding:"required"`
	ProblemText   string          `json:"problem_text" binding:"required"`
	Points        int32           `json:"points"`
	CorrectAnswer *string         `json:"correct_answer,omitempty"`
	Options       []OptionRequest `json:"options,omitempty"`
}

// UpdateQuizProblem updates a quiz problem
// PUT /api/quiz-problems/:id
func (h *Handler) UpdateQuizProblem(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid problem id"})
		return
	}

	var req UpdateQuizProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing problem
	existingProblem, err := h.queries.GetQuizProblem(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "quiz problem not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get quiz problem"})
		return
	}

	// Check institution access
	institution, err := h.getGroupInstitution(c, existingProblem.GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Set default points
	points := req.Points
	if points == 0 {
		points = 1
	}

	// Build correct answer field
	var correctAnswer pgtype.Text
	if req.CorrectAnswer != nil {
		correctAnswer = pgtype.Text{String: *req.CorrectAnswer, Valid: true}
	}

	problem, err := h.queries.UpdateQuizProblem(c.Request.Context(), sqlc.UpdateQuizProblemParams{
		ID:            id,
		GroupID:       existingProblem.GroupID,
		ProblemType:   sqlc.QuizProblemType(req.ProblemType),
		Name:          req.Name,
		Description:   req.Description,
		ProblemText:   req.ProblemText,
		Points:        points,
		CorrectAnswer: correctAnswer,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update quiz problem"})
		return
	}

	// Replace options: delete existing and create new
	if err := h.queries.DeleteQuizProblemOptionsByProblem(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update options"})
		return
	}

	var options []sqlc.QuizProblemOption
	for _, opt := range req.Options {
		option, err := h.queries.CreateQuizProblemOption(c.Request.Context(), sqlc.CreateQuizProblemOptionParams{
			ProblemID:    id,
			OptionText:   opt.OptionText,
			IsCorrect:    opt.IsCorrect,
			DisplayOrder: opt.DisplayOrder,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create option"})
			return
		}
		options = append(options, option)
	}

	c.JSON(http.StatusOK, QuizProblemResponse{
		Problem: problem,
		Options: options,
	})
}

// DeleteQuizProblem deletes a quiz problem
// DELETE /api/quiz-problems/:id
func (h *Handler) DeleteQuizProblem(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid problem id"})
		return
	}

	// Get existing problem to check institution
	existingProblem, err := h.queries.GetQuizProblem(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "quiz problem not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get quiz problem"})
		return
	}

	// Check institution access
	institution, err := h.getGroupInstitution(c, existingProblem.GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Options will be deleted via CASCADE
	if err := h.queries.DeleteQuizProblem(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete quiz problem"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "quiz problem deleted"})
}

// RegisterRoutes registers quiz problem routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware *auth.Middleware) {
	// Routes under quiz groups
	groupProblems := rg.Group("/quiz-groups/:id/problems")
	{
		groupProblems.GET("", authMiddleware.AuthRequired(), h.ListQuizProblems)
		groupProblems.POST("", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.CreateQuizProblem)
	}

	// Direct routes
	problems := rg.Group("/quiz-problems")
	{
		problems.GET("/:id", authMiddleware.AuthRequired(), h.GetQuizProblem)
		problems.PUT("/:id", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.UpdateQuizProblem)
		problems.DELETE("/:id", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.DeleteQuizProblem)
	}
}
