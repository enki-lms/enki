package quizsubmissions

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/enki/daemon/internal/auth"
	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Handler handles quiz submission-related API routes
type Handler struct {
	queries    *sqlc.Queries
	middleware *auth.Middleware
}

// NewHandler creates a new quiz submission handler
func NewHandler(queries *sqlc.Queries, middleware *auth.Middleware) *Handler {
	return &Handler{
		queries:    queries,
		middleware: middleware,
	}
}

// getProblemInstitution gets the institution for a quiz problem
func (h *Handler) getProblemInstitution(c *gin.Context, problemID int64) (string, error) {
	problem, err := h.queries.GetQuizProblem(c.Request.Context(), problemID)
	if err != nil {
		return "", err
	}
	group, err := h.queries.GetQuizProblemGroup(c.Request.Context(), problem.GroupID)
	if err != nil {
		return "", err
	}
	course, err := h.queries.GetCourse(c.Request.Context(), group.CourseID)
	if err != nil {
		return "", err
	}
	return course.Institution, nil
}

// SubmitQuizAnswerRequest represents the request body for submitting a quiz answer
type SubmitQuizAnswerRequest struct {
	AnswerText      *string `json:"answer_text,omitempty"`      // For open_ended, fill_blank
	SelectedOptions []int64 `json:"selected_options,omitempty"` // For true_false, mcq_single, mcq_multi
}

// SubmitQuizAnswer submits an answer to a quiz problem
// POST /api/quiz-problems/:id/submit
func (h *Handler) SubmitQuizAnswer(c *gin.Context) {
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

	var req SubmitQuizAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the problem
	problem, err := h.queries.GetQuizProblem(c.Request.Context(), problemID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "problem not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get problem"})
		return
	}

	// Verify user has access to this problem
	institution, err := h.getProblemInstitution(c, problemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Grade the submission based on problem type
	var isCorrect pgtype.Bool
	var score int32
	maxScore := problem.Points

	switch problem.ProblemType {
	case sqlc.QuizProblemTypeOpenEnded:
		// Open-ended questions are not auto-graded
		isCorrect = pgtype.Bool{Valid: false}
		score = 0 // Will be graded later

	case sqlc.QuizProblemTypeFillBlank:
		// Fill-in-the-blank: compare answer (case-insensitive)
		if req.AnswerText != nil && problem.CorrectAnswer.Valid {
			// Support multiple acceptable answers separated by commas
			acceptableAnswers := strings.Split(problem.CorrectAnswer.String, ",")
			userAnswer := strings.TrimSpace(strings.ToLower(*req.AnswerText))
			correct := false
			for _, acceptable := range acceptableAnswers {
				if strings.TrimSpace(strings.ToLower(acceptable)) == userAnswer {
					correct = true
					break
				}
			}
			isCorrect = pgtype.Bool{Bool: correct, Valid: true}
			if correct {
				score = maxScore
			}
		} else {
			isCorrect = pgtype.Bool{Bool: false, Valid: true}
		}

	case sqlc.QuizProblemTypeTrueFalse, sqlc.QuizProblemTypeMcqSingle, sqlc.QuizProblemTypeMcqMulti:
		// MCQ/True-False: compare selected options with correct options
		options, err := h.queries.ListQuizProblemOptionsByProblem(c.Request.Context(), problemID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get options"})
			return
		}

		// Find correct option IDs
		var correctOptionIDs []int64
		for _, opt := range options {
			if opt.IsCorrect {
				correctOptionIDs = append(correctOptionIDs, opt.ID)
			}
		}

		// Compare
		correct := compareOptions(req.SelectedOptions, correctOptionIDs)
		isCorrect = pgtype.Bool{Bool: correct, Valid: true}
		if correct {
			score = maxScore
		}
	}

	// Build answer_text for storage
	var answerText pgtype.Text
	if req.AnswerText != nil {
		answerText = pgtype.Text{String: *req.AnswerText, Valid: true}
	}

	// Create submission
	submission, err := h.queries.CreateQuizSubmission(c.Request.Context(), sqlc.CreateQuizSubmissionParams{
		UserID:          claims.UserID,
		ProblemID:       problemID,
		AnswerText:      answerText,
		SelectedOptions: req.SelectedOptions,
		IsCorrect:       isCorrect,
		Score:           score,
		MaxScore:        maxScore,
		Feedback:        pgtype.Text{Valid: false},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save submission"})
		return
	}

	c.JSON(http.StatusOK, submission)
}

// compareOptions checks if two slices of option IDs contain the same elements
func compareOptions(selected, correct []int64) bool {
	if len(selected) != len(correct) {
		return false
	}
	slices.Sort(selected)
	slices.Sort(correct)
	return slices.Equal(selected, correct)
}

// ListSubmissions returns submission history for a quiz problem
// GET /api/quiz-problems/:id/submissions
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
	submissions, err := h.queries.ListQuizSubmissionsByUserAndProblem(c.Request.Context(), sqlc.ListQuizSubmissionsByUserAndProblemParams{
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
// GET /api/quiz-submissions/:id
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

	submission, err := h.queries.GetQuizSubmission(c.Request.Context(), submissionID)
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

// GradeSubmissionRequest represents the request body for grading a submission
type GradeSubmissionRequest struct {
	IsCorrect *bool  `json:"is_correct"`
	Score     int32  `json:"score"`
	Feedback  string `json:"feedback"`
}

// GradeSubmission allows teachers to grade open-ended submissions
// POST /api/quiz-submissions/:id/grade
func (h *Handler) GradeSubmission(c *gin.Context) {
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

	var req GradeSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing submission
	existingSubmission, err := h.queries.GetQuizSubmission(c.Request.Context(), submissionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "submission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get submission"})
		return
	}

	// Verify teacher has access
	institution, err := h.getProblemInstitution(c, existingSubmission.ProblemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access"})
		return
	}

	if !h.middleware.CanAccessInstitution(claims, institution) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Build isCorrect
	var isCorrect pgtype.Bool
	if req.IsCorrect != nil {
		isCorrect = pgtype.Bool{Bool: *req.IsCorrect, Valid: true}
	}

	// Build feedback
	var feedback pgtype.Text
	if req.Feedback != "" {
		feedback = pgtype.Text{String: req.Feedback, Valid: true}
	}

	submission, err := h.queries.UpdateQuizSubmissionFeedback(c.Request.Context(), sqlc.UpdateQuizSubmissionFeedbackParams{
		ID:        submissionID,
		IsCorrect: isCorrect,
		Score:     req.Score,
		Feedback:  feedback,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update submission"})
		return
	}

	c.JSON(http.StatusOK, submission)
}

// ListAllSubmissionsForProblem returns all submissions for a problem (teachers only)
// GET /api/quiz-problems/:id/all-submissions
func (h *Handler) ListAllSubmissionsForProblem(c *gin.Context) {
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

	// Verify teacher has access
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

	submissions, err := h.queries.ListQuizSubmissionsByProblem(c.Request.Context(), problemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get submissions"})
		return
	}

	c.JSON(http.StatusOK, submissions)
}

// ListAllSubmissions returns all quiz submissions for the current user
// GET /api/quiz-submissions
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

	submissions, err := h.queries.ListQuizSubmissionsByUserWithLimit(c.Request.Context(), sqlc.ListQuizSubmissionsByUserWithLimitParams{
		UserID:     claims.UserID,
		LimitCount: limitCount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get submissions"})
		return
	}

	c.JSON(http.StatusOK, submissions)
}

// RegisterRoutes registers quiz submission routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware *auth.Middleware) {
	// Problem submission routes
	problemSubmissions := rg.Group("/quiz-problems/:id")
	problemSubmissions.Use(authMiddleware.AuthRequired())
	{
		problemSubmissions.POST("/submit", h.SubmitQuizAnswer)
		problemSubmissions.GET("/submissions", h.ListSubmissions)
		problemSubmissions.GET("/all-submissions", authMiddleware.RequireTeacher(), h.ListAllSubmissionsForProblem)
	}

	// Direct submission routes
	submissions := rg.Group("/quiz-submissions")
	submissions.Use(authMiddleware.AuthRequired())
	{
		submissions.GET("", h.ListAllSubmissions)
		submissions.GET("/:id", h.GetSubmission)
		submissions.POST("/:id/grade", authMiddleware.RequireTeacher(), h.GradeSubmission)
	}
}
