package quizgroups

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/enki/daemon/internal/auth"
	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// Handler handles quiz problem group-related API routes
type Handler struct {
	queries    *sqlc.Queries
	middleware *auth.Middleware
}

// NewHandler creates a new quiz problem group handler
func NewHandler(queries *sqlc.Queries, middleware *auth.Middleware) *Handler {
	return &Handler{
		queries:    queries,
		middleware: middleware,
	}
}

// getCourseInstitution gets the institution for a course
func (h *Handler) getCourseInstitution(c *gin.Context, courseID int64) (string, error) {
	course, err := h.queries.GetCourse(c.Request.Context(), courseID)
	if err != nil {
		return "", err
	}
	return course.Institution, nil
}

// getGroupInstitution gets the institution for a quiz problem group via its course
func (h *Handler) getGroupInstitution(c *gin.Context, groupID int64) (string, error) {
	group, err := h.queries.GetQuizProblemGroup(c.Request.Context(), groupID)
	if err != nil {
		return "", err
	}
	return h.getCourseInstitution(c, group.CourseID)
}

// ListQuizProblemGroups returns quiz problem groups for a course
// GET /api/courses/:id/quiz-groups
func (h *Handler) ListQuizProblemGroups(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	// Verify course belongs to user's institution
	institution, err := h.getCourseInstitution(c, courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get course"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	groups, err := h.queries.ListQuizProblemGroupsByCourse(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list quiz problem groups"})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// GetQuizProblemGroup returns a quiz problem group by ID
// GET /api/quiz-groups/:id
func (h *Handler) GetQuizProblemGroup(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	group, err := h.queries.GetQuizProblemGroup(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "quiz problem group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get quiz problem group"})
		return
	}

	// Check institution access via course
	institution, err := h.getCourseInstitution(c, group.CourseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// CreateQuizProblemGroupRequest represents the request body for creating a quiz problem group
type CreateQuizProblemGroupRequest struct {
	Type        string `json:"type" binding:"required,oneof=exam practice"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// CreateQuizProblemGroup creates a new quiz problem group
// POST /api/courses/:id/quiz-groups
func (h *Handler) CreateQuizProblemGroup(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	var req CreateQuizProblemGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify course belongs to user's institution
	institution, err := h.getCourseInstitution(c, courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get course"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	group, err := h.queries.CreateQuizProblemGroup(c.Request.Context(), sqlc.CreateQuizProblemGroupParams{
		Type:        sqlc.CompSciProblemType(req.Type),
		CourseID:    courseID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create quiz problem group"})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// UpdateQuizProblemGroupRequest represents the request body for updating a quiz problem group
type UpdateQuizProblemGroupRequest struct {
	Type        string `json:"type" binding:"required,oneof=exam practice"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// UpdateQuizProblemGroup updates a quiz problem group
// PUT /api/quiz-groups/:id
func (h *Handler) UpdateQuizProblemGroup(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	var req UpdateQuizProblemGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing group
	existingGroup, err := h.queries.GetQuizProblemGroup(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "quiz problem group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get quiz problem group"})
		return
	}

	// Check institution access via course
	institution, err := h.getCourseInstitution(c, existingGroup.CourseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	group, err := h.queries.UpdateQuizProblemGroup(c.Request.Context(), sqlc.UpdateQuizProblemGroupParams{
		ID:          id,
		Type:        sqlc.CompSciProblemType(req.Type),
		CourseID:    existingGroup.CourseID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update quiz problem group"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// DeleteQuizProblemGroup deletes a quiz problem group
// DELETE /api/quiz-groups/:id
func (h *Handler) DeleteQuizProblemGroup(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	// Get existing group to check institution
	existingGroup, err := h.queries.GetQuizProblemGroup(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "quiz problem group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get quiz problem group"})
		return
	}

	// Check institution access via course
	institution, err := h.getCourseInstitution(c, existingGroup.CourseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := h.queries.DeleteQuizProblemGroup(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete quiz problem group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "quiz problem group deleted"})
}

// RegisterRoutes registers quiz problem group routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware *auth.Middleware) {
	// Routes under courses
	courseGroups := rg.Group("/courses/:id/quiz-groups")
	{
		courseGroups.GET("", authMiddleware.AuthRequired(), h.ListQuizProblemGroups)
		courseGroups.POST("", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.CreateQuizProblemGroup)
	}

	// Direct routes
	groups := rg.Group("/quiz-groups")
	{
		groups.GET("/:id", authMiddleware.AuthRequired(), h.GetQuizProblemGroup)
		groups.PUT("/:id", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.UpdateQuizProblemGroup)
		groups.DELETE("/:id", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.DeleteQuizProblemGroup)
	}
}
