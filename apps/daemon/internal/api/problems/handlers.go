package problems

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

// Handler handles problem-related API routes
type Handler struct {
	queries    *sqlc.Queries
	middleware *auth.Middleware
}

// NewHandler creates a new problem handler
func NewHandler(queries *sqlc.Queries, middleware *auth.Middleware) *Handler {
	return &Handler{
		queries:    queries,
		middleware: middleware,
	}
}

// getProblemInstitution gets the institution for a problem via group -> course
func (h *Handler) getProblemInstitution(c *gin.Context, problemID int64) (string, error) {
	problem, err := h.queries.GetCompSciProblem(c.Request.Context(), problemID)
	if err != nil {
		return "", err
	}
	return h.getGroupInstitution(c, problem.GroupID)
}

// getGroupInstitution gets the institution for a group via course
func (h *Handler) getGroupInstitution(c *gin.Context, groupID int64) (string, error) {
	group, err := h.queries.GetCompSciProblemGroup(c.Request.Context(), groupID)
	if err != nil {
		return "", err
	}
	course, err := h.queries.GetCourse(c.Request.Context(), group.CourseID)
	if err != nil {
		return "", err
	}
	return course.Institution, nil
}

// ListProblems returns problems for a group
// GET /api/problem-groups/:id/problems
func (h *Handler) ListProblems(c *gin.Context) {
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

	// Verify group's course belongs to user's institution
	institution, err := h.getGroupInstitution(c, groupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "problem group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	problems, err := h.queries.ListCompSciProblemsByGroup(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list problems"})
		return
	}

	c.JSON(http.StatusOK, problems)
}

// GetProblem returns a problem by ID
// GET /api/problems/:id
func (h *Handler) GetProblem(c *gin.Context) {
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

	problem, err := h.queries.GetCompSciProblem(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "problem not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get problem"})
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

	c.JSON(http.StatusOK, problem)
}

// CreateProblemRequest represents the request body for creating a problem
type CreateProblemRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description" binding:"required"`
	ProblemText   string `json:"problem_text" binding:"required"`
	TimeLimitMs   *int32 `json:"time_limit_ms"`
	MemoryLimitMB *int32 `json:"memory_limit_mb"`
}

// CreateProblem creates a new problem
// POST /api/problem-groups/:id/problems
func (h *Handler) CreateProblem(c *gin.Context) {
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

	var req CreateProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify group's course belongs to user's institution
	institution, err := h.getGroupInstitution(c, groupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "problem group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Convert optional limits to pgtype
	var timeLimitMs, memoryLimitMb pgtype.Int4
	if req.TimeLimitMs != nil {
		timeLimitMs = pgtype.Int4{Int32: *req.TimeLimitMs, Valid: true}
	}
	if req.MemoryLimitMB != nil {
		memoryLimitMb = pgtype.Int4{Int32: *req.MemoryLimitMB, Valid: true}
	}

	problem, err := h.queries.CreateCompSciProblem(c.Request.Context(), sqlc.CreateCompSciProblemParams{
		GroupID:       groupID,
		Name:          req.Name,
		Description:   req.Description,
		ProblemText:   req.ProblemText,
		TimeLimitMs:   timeLimitMs,
		MemoryLimitMb: memoryLimitMb,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create problem"})
		return
	}

	c.JSON(http.StatusCreated, problem)
}

// UpdateProblemRequest represents the request body for updating a problem
type UpdateProblemRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description" binding:"required"`
	ProblemText   string `json:"problem_text" binding:"required"`
	TimeLimitMs   *int32 `json:"time_limit_ms"`
	MemoryLimitMB *int32 `json:"memory_limit_mb"`
}

// UpdateProblem updates a problem
// PUT /api/problems/:id
func (h *Handler) UpdateProblem(c *gin.Context) {
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

	var req UpdateProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing problem
	existingProblem, err := h.queries.GetCompSciProblem(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "problem not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get problem"})
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

	// Convert optional limits to pgtype
	var timeLimitMs, memoryLimitMb pgtype.Int4
	if req.TimeLimitMs != nil {
		timeLimitMs = pgtype.Int4{Int32: *req.TimeLimitMs, Valid: true}
	}
	if req.MemoryLimitMB != nil {
		memoryLimitMb = pgtype.Int4{Int32: *req.MemoryLimitMB, Valid: true}
	}

	problem, err := h.queries.UpdateCompSciProblem(c.Request.Context(), sqlc.UpdateCompSciProblemParams{
		ID:            id,
		GroupID:       existingProblem.GroupID,
		Name:          req.Name,
		Description:   req.Description,
		ProblemText:   req.ProblemText,
		TimeLimitMs:   timeLimitMs,
		MemoryLimitMb: memoryLimitMb,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update problem"})
		return
	}

	c.JSON(http.StatusOK, problem)
}

// DeleteProblem deletes a problem
// DELETE /api/problems/:id
func (h *Handler) DeleteProblem(c *gin.Context) {
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
	existingProblem, err := h.queries.GetCompSciProblem(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "problem not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get problem"})
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

	if err := h.queries.DeleteCompSciProblem(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete problem"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "problem deleted"})
}

// RegisterRoutes registers problem routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware *auth.Middleware) {
	// Routes under groups
	groupProblems := rg.Group("/problem-groups/:id/problems")
	{
		groupProblems.GET("", authMiddleware.AuthRequired(), h.ListProblems)
		groupProblems.POST("", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.CreateProblem)
	}

	// Direct routes
	problems := rg.Group("/problems")
	{
		problems.GET("/:id", authMiddleware.AuthRequired(), h.GetProblem)
		problems.PUT("/:id", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.UpdateProblem)
		problems.DELETE("/:id", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.DeleteProblem)
	}
}
