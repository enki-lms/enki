package courses

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/enki/daemon/internal/auth"
	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// Handler handles course-related API routes
type Handler struct {
	queries    *sqlc.Queries
	middleware *auth.Middleware
}

// NewHandler creates a new course handler
func NewHandler(queries *sqlc.Queries, middleware *auth.Middleware) *Handler {
	return &Handler{
		queries:    queries,
		middleware: middleware,
	}
}

// ListCourses returns courses from the user's institution
// GET /api/courses
func (h *Handler) ListCourses(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courses, err := h.queries.ListCoursesByUser(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// GetCourse returns a course by ID
// GET /api/courses/:id
func (h *Handler) GetCourse(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	course, err := h.queries.GetCourse(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get course"})
		return
	}

	// Check institution access
	if !h.middleware.CanAccessInstitution(claims, course.Institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, course)
}

// CreateCourseRequest represents the request body for creating a course
type CreateCourseRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateCourse creates a new course
// POST /api/courses
func (h *Handler) CreateCourse(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Auto-set institution from the teacher's institution
	course, err := h.queries.CreateCourse(c.Request.Context(), sqlc.CreateCourseParams{
		Name:        req.Name,
		Institution: claims.Institution,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// UpdateCourseRequest represents the request body for updating a course
type UpdateCourseRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateCourse updates a course
// PUT /api/courses/:id
func (h *Handler) UpdateCourse(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	var req UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing course to check institution
	existingCourse, err := h.queries.GetCourse(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get course"})
		return
	}

	// Check institution access
	if !h.middleware.CanAccessInstitution(claims, existingCourse.Institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	course, err := h.queries.UpdateCourse(c.Request.Context(), sqlc.UpdateCourseParams{
		ID:          id,
		Name:        req.Name,
		Institution: existingCourse.Institution,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update course"})
		return
	}

	c.JSON(http.StatusOK, course)
}

// DeleteCourse deletes a course
// DELETE /api/courses/:id
func (h *Handler) DeleteCourse(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	// Get existing course to check institution
	existingCourse, err := h.queries.GetCourse(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get course"})
		return
	}

	// Check institution access
	if !h.middleware.CanAccessInstitution(claims, existingCourse.Institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := h.queries.DeleteCourse(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "course deleted"})
}

// RegisterRoutes registers course routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware *auth.Middleware) {
	courses := rg.Group("/courses")
	{
		// Read operations - authenticated users
		courses.GET("", authMiddleware.AuthRequired(), h.ListCourses)
		courses.GET("/:id", authMiddleware.AuthRequired(), h.GetCourse)

		// Write operations - teachers only
		courses.POST("", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.CreateCourse)
		courses.PUT("/:id", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.UpdateCourse)
		courses.DELETE("/:id", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.DeleteCourse)
	}
}
