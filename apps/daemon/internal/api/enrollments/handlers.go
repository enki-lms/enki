package enrollments

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/enki/daemon/internal/auth"
	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// Handler handles enrollment-related API routes
type Handler struct {
	queries    *sqlc.Queries
	middleware *auth.Middleware
}

// NewHandler creates a new enrollment handler
func NewHandler(queries *sqlc.Queries, middleware *auth.Middleware) *Handler {
	return &Handler{
		queries:    queries,
		middleware: middleware,
	}
}

// ListEnrollments returns enrolled users for a course
// GET /api/courses/:id/enrollments
func (h *Handler) ListEnrollments(c *gin.Context) {
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
	course, err := h.queries.GetCourse(c.Request.Context(), courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get course"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, course.Institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	enrollments, err := h.queries.ListUserCoursesByCourse(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list enrollments"})
		return
	}

	c.JSON(http.StatusOK, enrollments)
}

// EnrollUserRequest represents the request body for enrolling a user
type EnrollUserRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

// EnrollUser enrolls a user in a course
// POST /api/courses/:id/enrollments
func (h *Handler) EnrollUser(c *gin.Context) {
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

	var req EnrollUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify course belongs to user's institution
	course, err := h.queries.GetCourse(c.Request.Context(), courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get course"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, course.Institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Verify user belongs to same institution
	user, err := h.queries.GetUser(c.Request.Context(), req.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	if user.Institution != course.Institution {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user and course must be in same institution"})
		return
	}

	enrollment, err := h.queries.CreateUserCourse(c.Request.Context(), sqlc.CreateUserCourseParams{
		UserID:   req.UserID,
		CourseID: courseID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enroll user"})
		return
	}

	c.JSON(http.StatusCreated, enrollment)
}

// UnenrollUser removes a user from a course
// DELETE /api/courses/:id/enrollments/:userId
func (h *Handler) UnenrollUser(c *gin.Context) {
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

	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	// Verify course belongs to user's institution
	course, err := h.queries.GetCourse(c.Request.Context(), courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get course"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, course.Institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := h.queries.DeleteUserCourseByUserAndCourse(c.Request.Context(), sqlc.DeleteUserCourseByUserAndCourseParams{
		UserID:   userID,
		CourseID: courseID,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unenroll user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user unenrolled"})
}

// RegisterRoutes registers enrollment routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware *auth.Middleware) {
	enrollments := rg.Group("/courses/:id/enrollments")
	enrollments.Use(authMiddleware.AuthRequired(), authMiddleware.RequireTeacher())
	{
		enrollments.GET("", h.ListEnrollments)
		enrollments.POST("", h.EnrollUser)
		enrollments.DELETE("/:userId", h.UnenrollUser)
	}
}
