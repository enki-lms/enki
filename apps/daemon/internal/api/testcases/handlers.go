package testcases

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/enki/daemon/internal/auth"
	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// Handler handles test case-related API routes
type Handler struct {
	queries    *sqlc.Queries
	middleware *auth.Middleware
}

// NewHandler creates a new test case handler
func NewHandler(queries *sqlc.Queries, middleware *auth.Middleware) *Handler {
	return &Handler{
		queries:    queries,
		middleware: middleware,
	}
}

// getTestCaseInstitution gets the institution for a test case via problem -> group -> course
func (h *Handler) getTestCaseInstitution(c *gin.Context, testCaseID int64) (string, error) {
	testCase, err := h.queries.GetCompSciTestCase(c.Request.Context(), testCaseID)
	if err != nil {
		return "", err
	}
	return h.getProblemInstitution(c, testCase.ProblemID)
}

// getProblemInstitution gets the institution for a problem via group -> course
func (h *Handler) getProblemInstitution(c *gin.Context, problemID int64) (string, error) {
	problem, err := h.queries.GetCompSciProblem(c.Request.Context(), problemID)
	if err != nil {
		return "", err
	}
	group, err := h.queries.GetCompSciProblemGroup(c.Request.Context(), problem.GroupID)
	if err != nil {
		return "", err
	}
	course, err := h.queries.GetCourse(c.Request.Context(), group.CourseID)
	if err != nil {
		return "", err
	}
	return course.Institution, nil
}

// ListTestCases returns test cases for a problem
// GET /api/problems/:id/test-cases
func (h *Handler) ListTestCases(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	problemID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid problem id"})
		return
	}

	// Verify problem's institution
	institution, err := h.getProblemInstitution(c, problemID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "problem not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	testCases, err := h.queries.ListCompSciTestCasesByProblem(c.Request.Context(), problemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list test cases"})
		return
	}

	c.JSON(http.StatusOK, testCases)
}

// GetTestCase returns a test case by ID
// GET /api/test-cases/:id
func (h *Handler) GetTestCase(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test case id"})
		return
	}

	testCase, err := h.queries.GetCompSciTestCase(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "test case not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get test case"})
		return
	}

	// Check institution access
	institution, err := h.getProblemInstitution(c, testCase.ProblemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, testCase)
}

// CreateTestCaseRequest represents the request body for creating a test case
type CreateTestCaseRequest struct {
	Input         string `json:"input" binding:"required"`
	Output        string `json:"output" binding:"required"`
	CorrectPoints int32  `json:"correct_points" binding:"required"`
}

// CreateTestCase creates a new test case
// POST /api/problems/:id/test-cases
func (h *Handler) CreateTestCase(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	problemID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid problem id"})
		return
	}

	var req CreateTestCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify problem's institution
	institution, err := h.getProblemInstitution(c, problemID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "problem not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	testCase, err := h.queries.CreateCompSciTestCase(c.Request.Context(), sqlc.CreateCompSciTestCaseParams{
		ProblemID:     problemID,
		Input:         req.Input,
		Output:        req.Output,
		CorrectPoints: req.CorrectPoints,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create test case"})
		return
	}

	c.JSON(http.StatusCreated, testCase)
}

// UpdateTestCaseRequest represents the request body for updating a test case
type UpdateTestCaseRequest struct {
	Input         string `json:"input" binding:"required"`
	Output        string `json:"output" binding:"required"`
	CorrectPoints int32  `json:"correct_points" binding:"required"`
}

// UpdateTestCase updates a test case
// PUT /api/test-cases/:id
func (h *Handler) UpdateTestCase(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test case id"})
		return
	}

	var req UpdateTestCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing test case
	existingTestCase, err := h.queries.GetCompSciTestCase(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "test case not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get test case"})
		return
	}

	// Check institution access
	institution, err := h.getProblemInstitution(c, existingTestCase.ProblemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	testCase, err := h.queries.UpdateCompSciTestCase(c.Request.Context(), sqlc.UpdateCompSciTestCaseParams{
		ID:            id,
		ProblemID:     existingTestCase.ProblemID,
		Input:         req.Input,
		Output:        req.Output,
		CorrectPoints: req.CorrectPoints,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update test case"})
		return
	}

	c.JSON(http.StatusOK, testCase)
}

// DeleteTestCase deletes a test case
// DELETE /api/test-cases/:id
func (h *Handler) DeleteTestCase(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test case id"})
		return
	}

	// Get existing test case to check institution
	existingTestCase, err := h.queries.GetCompSciTestCase(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "test case not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get test case"})
		return
	}

	// Check institution access
	institution, err := h.getProblemInstitution(c, existingTestCase.ProblemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := h.queries.DeleteCompSciTestCase(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete test case"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "test case deleted"})
}

// RegisterRoutes registers test case routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware *auth.Middleware) {
	// Routes under problems
	problemTestCases := rg.Group("/problems/:id/test-cases")
	{
		problemTestCases.GET("", authMiddleware.AuthRequired(), h.ListTestCases)
		problemTestCases.POST("", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.CreateTestCase)
	}

	// Direct routes
	testCases := rg.Group("/test-cases")
	{
		testCases.GET("/:id", authMiddleware.AuthRequired(), h.GetTestCase)
		testCases.PUT("/:id", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.UpdateTestCase)
		testCases.DELETE("/:id", authMiddleware.AuthRequired(), authMiddleware.RequireTeacher(), h.DeleteTestCase)
	}
}
