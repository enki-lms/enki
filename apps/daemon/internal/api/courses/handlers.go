package courses

import (
	"errors"
	"net/http"
	"sort"
	"strconv"

	"github.com/enki/daemon/internal/auth"
	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

// ListTeacherCourses returns courses owned by the authenticated teacher
// GET /api/courses/teaching
func (h *Handler) ListTeacherCourses(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courses, err := h.queries.ListCoursesByOwner(c.Request.Context(), pgtype.Int8{Int64: claims.UserID, Valid: true})
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
	Name    string `json:"name" binding:"required"`
	Subject string `json:"subject" binding:"required"`
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
		Subject:     req.Subject,
		Institution: claims.Institution,
		OwnerID:     pgtype.Int8{Int64: claims.UserID, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// UpdateCourseRequest represents the request body for updating a course
type UpdateCourseRequest struct {
	Name    string `json:"name" binding:"required"`
	Subject string `json:"subject" binding:"required"`
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
		Subject:     req.Subject,
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

// ListStudentSubmissions returns submission history for a specific student in a course
// GET /api/courses/:id/students/:studentId/submissions
func (h *Handler) ListStudentSubmissions(c *gin.Context) {
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

	studentID, err := strconv.ParseInt(c.Param("studentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	// 1. Verify access: User must be teacher of the course (or admin)
	course, err := h.queries.GetCourse(c.Request.Context(), courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get course"})
		return
	}

	// Only allow the owner of the course or an admin to view student submissions
	isOwner := course.OwnerID.Valid && course.OwnerID.Int64 == claims.UserID
	isAdmin := claims.Role == "admin"

	if !isOwner && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// 2. Fetch CompSci submissions
	compSciSubmissions, err := h.queries.ListCompSciSubmissionsByCourseAndUser(c.Request.Context(), sqlc.ListCompSciSubmissionsByCourseAndUserParams{
		CourseID: courseID,
		UserID:   studentID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch computer science submissions"})
		return
	}

	// 3. Fetch Quiz submissions
	quizSubmissions, err := h.queries.ListQuizSubmissionsByCourseAndUser(c.Request.Context(), sqlc.ListQuizSubmissionsByCourseAndUserParams{
		CourseID: courseID,
		UserID:   studentID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch quiz submissions"})
		return
	}

	// 4. Combine and sort
	type SubmissionHistoryItem struct {
		Type      string           `json:"type"` // "comp_sci" or "quiz"
		ID        int64            `json:"id"`
		CreatedAt pgtype.Timestamp `json:"created_at"`
		UserID    int64            `json:"user_id"`
		ProblemID int64            `json:"problem_id"`
		Score     int32            `json:"score"`
		MaxScore  int32            `json:"max_score"`
		// CompSci specific
		Code        string `json:"code,omitempty"`
		PassedTests int32  `json:"passed_tests,omitempty"`
		TotalTests  int32  `json:"total_tests,omitempty"`
		ResultsJson string `json:"results_json,omitempty"`
		// Quiz specific
		AnswerText      string  `json:"answer_text,omitempty"`
		SelectedOptions []int64 `json:"selected_options,omitempty"`
		IsCorrect       *bool   `json:"is_correct,omitempty"` // Pointer to handle null/false distinction if needed
		Feedback        string  `json:"feedback,omitempty"`
	}

	var history []SubmissionHistoryItem

	for _, s := range compSciSubmissions {
		history = append(history, SubmissionHistoryItem{
			Type:        "comp_sci",
			ID:          s.ID,
			CreatedAt:   s.CreatedAt,
			UserID:      s.UserID,
			ProblemID:   s.ProblemID,
			Score:       s.Score,
			MaxScore:    s.MaxScore,
			Code:        s.Code,
			PassedTests: s.PassedTests,
			TotalTests:  s.TotalTests,
			ResultsJson: s.ResultsJson,
		})
	}

	for _, s := range quizSubmissions {
		var isCorrect *bool
		if s.IsCorrect.Valid {
			val := s.IsCorrect.Bool
			isCorrect = &val
		}

		var answerText string
		if s.AnswerText.Valid {
			answerText = s.AnswerText.String
		}

		var feedback string
		if s.Feedback.Valid {
			feedback = s.Feedback.String
		}

		history = append(history, SubmissionHistoryItem{
			Type:            "quiz",
			ID:              s.ID,
			CreatedAt:       s.CreatedAt,
			UserID:          s.UserID,
			ProblemID:       s.ProblemID,
			Score:           s.Score,
			MaxScore:        s.MaxScore,
			AnswerText:      answerText,
			SelectedOptions: s.SelectedOptions,
			IsCorrect:       isCorrect,
			Feedback:        feedback,
		})
	}

	// Sort by CreatedAt descending
	sort.Slice(history, func(i, j int) bool {
		return history[i].CreatedAt.Time.After(history[j].CreatedAt.Time)
	})

	c.JSON(http.StatusOK, history)
}

// RegisterRoutes registers course routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware *auth.Middleware) {
	courses := rg.Group("/courses")
	{
		// Read operations - authenticated users
		courses.GET("", authMiddleware.AuthRequired(), h.ListCourses)
		courses.GET("/teaching", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.ListTeacherCourses)
		courses.GET("/:id", authMiddleware.AuthRequired(), h.GetCourse)
		courses.GET("/:id/students/:studentId/submissions", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.ListStudentSubmissions)

		// Write operations - teachers only
		courses.POST("", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.CreateCourse)
		courses.PUT("/:id", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.UpdateCourse)
		courses.DELETE("/:id", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.DeleteCourse)
	}
}
