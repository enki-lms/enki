package submissions

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/enki/daemon/internal/auth"
	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/enki/daemon/internal/grading"
	"github.com/enki/daemon/internal/problem_eval"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Handler handles submission-related API routes
type Handler struct {
	queries        *sqlc.Queries
	middleware     *auth.Middleware
	gradingService *grading.Service
}

// NewHandler creates a new submission handler
func NewHandler(queries *sqlc.Queries, middleware *auth.Middleware, gradingService *grading.Service) *Handler {
	return &Handler{
		queries:        queries,
		middleware:     middleware,
		gradingService: gradingService,
	}
}

// SubmitCodeRequest represents the request body for submitting code
type SubmitCodeRequest struct {
	Code string `json:"code" binding:"required"`
}

// SubmitCode runs code against all test cases for a problem
// POST /api/problems/:problemId/submit
func (h *Handler) SubmitCode(c *gin.Context) {
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

	var req SubmitCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the problem
	problem, err := h.queries.GetCompSciProblem(c.Request.Context(), problemID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "problem not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get problem"})
		return
	}

	// Verify user has access to this problem via group -> course -> institution
	institution, err := h.getProblemInstitution(c.Request.Context(), problem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Get all test cases for this problem
	testCases, err := h.queries.ListCompSciTestCasesByProblem(c.Request.Context(), problemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get test cases"})
		return
	}

	if len(testCases) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no test cases defined for this problem"})
		return
	}

	// Get per-problem limits (or use defaults via 0)
	var timeoutMs, memoryLimitMB int
	if problem.TimeLimitMs.Valid {
		timeoutMs = int(problem.TimeLimitMs.Int32)
	}
	if problem.MemoryLimitMb.Valid {
		memoryLimitMB = int(problem.MemoryLimitMb.Int32)
	}

	// Run code against each test case using grading service
	gradingResult, err := h.gradingService.GradeCompSci(c.Request.Context(), req.Code, testCases, timeoutMs, memoryLimitMB, problem.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to grade submission"})
		return
	}

	// Serialize results to JSON for storage
	resultsJSON, err := json.Marshal(gradingResult.Results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize results"})
		return
	}

	// Save submission to database
	submission, err := h.queries.CreateCompSciSubmission(c.Request.Context(), sqlc.CreateCompSciSubmissionParams{
		UserID:      claims.UserID,
		ProblemID:   problemID,
		Code:        req.Code,
		Score:       gradingResult.Score,
		MaxScore:    gradingResult.MaxScore,
		PassedTests: gradingResult.PassedTests,
		TotalTests:  gradingResult.TotalTests,
		ResultsJson: string(resultsJSON),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save submission"})
		return
	}

	response := problem_eval.SubmissionResult{
		ProblemID:      problemID,
		TotalTestCases: int(gradingResult.TotalTests),
		Passed:         int(gradingResult.PassedTests),
		Failed:         int(gradingResult.TotalTests - gradingResult.PassedTests),
		Score:          gradingResult.Score,
		MaxScore:       gradingResult.MaxScore,
		Results:        gradingResult.Results,
	}

	c.JSON(http.StatusOK, gin.H{
		"submission_id": submission.ID,
		"result":        response,
	})
}

// ListSubmissions returns submission history for a problem
// GET /api/problems/:problemId/submissions
func (h *Handler) ListSubmissions(c *gin.Context) {
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

	// Get submissions for this user and problem
	submissions, err := h.queries.ListCompSciSubmissionsByUserAndProblem(c.Request.Context(), sqlc.ListCompSciSubmissionsByUserAndProblemParams{
		UserID:    claims.UserID,
		ProblemID: problemID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get submissions"})
		return
	}

	c.JSON(http.StatusOK, submissions)
}

// GetSubmission returns a single submission
// GET /api/submissions/:id
func (h *Handler) GetSubmission(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	submissionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid submission id"})
		return
	}

	submission, err := h.queries.GetCompSciSubmission(c.Request.Context(), submissionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "submission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get submission"})
		return
	}

	// Only allow user to see their own submission (or teacher to see any)
	if submission.UserID != claims.UserID && !h.middleware.IsTeacherOrAdmin(claims) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, submission)
}

// getProblemInstitution gets the institution for a problem via group -> course
func (h *Handler) getProblemInstitution(ctx context.Context, problem sqlc.CompSciProblem) (string, error) {
	group, err := h.queries.GetCompSciProblemGroup(ctx, problem.GroupID)
	if err != nil {
		return "", err
	}
	course, err := h.queries.GetCourse(ctx, group.CourseID)
	if err != nil {
		return "", err
	}
	return course.Institution, nil
}

// ListAllSubmissions returns all submissions for the current user
// GET /api/submissions
func (h *Handler) ListAllSubmissions(c *gin.Context) {
	claims, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse optional limit parameter
	var limitCount pgtype.Int4
	if limitStr := c.Query("limit"); limitStr != "" {
		limitVal, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || limitVal <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
			return
		}
		limitCount = pgtype.Int4{Int32: int32(limitVal), Valid: true}
	}

	submissions, err := h.queries.ListCompSciSubmissionsByUserWithLimit(c.Request.Context(), sqlc.ListCompSciSubmissionsByUserWithLimitParams{
		UserID:     claims.UserID,
		LimitCount: limitCount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get submissions"})
		return
	}

	c.JSON(http.StatusOK, submissions)
}

// RegisterRoutes registers submission routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware *auth.Middleware) {
	problemSubmissions := rg.Group("/problems/:id")
	problemSubmissions.Use(authMiddleware.AuthRequired())
	{
		problemSubmissions.POST("/submit", h.SubmitCode)
		problemSubmissions.GET("/submissions", h.ListSubmissions)
	}

	submissions := rg.Group("/submissions")
	submissions.Use(authMiddleware.AuthRequired())
	{
		submissions.GET("", h.ListAllSubmissions)
		submissions.GET("/:id", h.GetSubmission)
	}
}
