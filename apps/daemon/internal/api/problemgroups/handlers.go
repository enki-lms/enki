package problemgroups

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/enki/daemon/internal/auth"
	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// Handler handles problem group-related API routes
type Handler struct {
	queries    *sqlc.Queries
	middleware *auth.Middleware
}

// NewHandler creates a new problem group handler
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

// getGroupInstitution gets the institution for a problem group via its course
func (h *Handler) getGroupInstitution(c *gin.Context, groupID int64) (string, error) {
	group, err := h.queries.GetCompSciProblemGroup(c.Request.Context(), groupID)
	if err != nil {
		return "", err
	}
	return h.getCourseInstitution(c, group.CourseID)
}

// ListProblemGroups returns problem groups for a course
// GET /api/courses/:id/problem-groups
func (h *Handler) ListProblemGroups(c *gin.Context) {
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

	groups, err := h.queries.ListCompSciProblemGroupsByCourse(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list problem groups"})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// GetProblemGroup returns a problem group by ID
// GET /api/problem-groups/:id
func (h *Handler) GetProblemGroup(c *gin.Context) {
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

	group, err := h.queries.GetCompSciProblemGroup(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "problem group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get problem group"})
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

// CreateProblemGroupRequest represents the request body for creating a problem group
type CreateProblemGroupRequest struct {
	Type        string `json:"type" binding:"required,oneof=exam practice"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// CreateProblemGroup creates a new problem group
// POST /api/courses/:id/problem-groups
func (h *Handler) CreateProblemGroup(c *gin.Context) {
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

	var req CreateProblemGroupRequest
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

	group, err := h.queries.CreateCompSciProblemGroup(c.Request.Context(), sqlc.CreateCompSciProblemGroupParams{
		Type:        sqlc.CompSciProblemType(req.Type),
		CourseID:    courseID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create problem group"})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// UpdateProblemGroupRequest represents the request body for updating a problem group
type UpdateProblemGroupRequest struct {
	Type        string `json:"type" binding:"required,oneof=exam practice"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// UpdateProblemGroup updates a problem group
// PUT /api/problem-groups/:id
func (h *Handler) UpdateProblemGroup(c *gin.Context) {
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

	var req UpdateProblemGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing group
	existingGroup, err := h.queries.GetCompSciProblemGroup(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "problem group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get problem group"})
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

	group, err := h.queries.UpdateCompSciProblemGroup(c.Request.Context(), sqlc.UpdateCompSciProblemGroupParams{
		ID:          id,
		Type:        sqlc.CompSciProblemType(req.Type),
		CourseID:    existingGroup.CourseID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update problem group"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// DeleteProblemGroup deletes a problem group
// DELETE /api/problem-groups/:id
func (h *Handler) DeleteProblemGroup(c *gin.Context) {
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
	existingGroup, err := h.queries.GetCompSciProblemGroup(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "problem group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get problem group"})
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

	if err := h.queries.DeleteCompSciProblemGroup(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete problem group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "problem group deleted"})
}

// RegisterRoutes registers problem group routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware *auth.Middleware) {
	// Routes under courses
	courseGroups := rg.Group("/courses/:id/problem-groups")
	{
		courseGroups.GET("", authMiddleware.AuthRequired(), h.ListProblemGroups)
		courseGroups.POST("", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.CreateProblemGroup)
	}

	// Direct routes
	groups := rg.Group("/problem-groups")
	{
		groups.GET("/:id", authMiddleware.AuthRequired(), h.GetProblemGroup)
		groups.PUT("/:id", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.UpdateProblemGroup)
		groups.DELETE("/:id", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.DeleteProblemGroup)
	}
}
