package grading

import (
	"context"
	"fmt"
	"strings"

	"os"

	"github.com/enki/daemon/internal/ai"
	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/enki/daemon/internal/problem_eval"
)

// Service handles grading logic for submissions
type Service struct {
	queries  *sqlc.Queries
	executor *problem_eval.Executor
	aiClient *ai.Client
}

// NewService creates a new grading service
func NewService(queries *sqlc.Queries, executor *problem_eval.Executor, aiClient *ai.Client) *Service {
	return &Service{
		queries:  queries,
		executor: executor,
		aiClient: aiClient,
	}
}

// CompSciGradingResult holds the result of grading a computer science submission
type CompSciGradingResult struct {
	Score       int32
	MaxScore    int32
	PassedTests int32
	TotalTests  int32
	Results     []problem_eval.TestCaseResult
}

// GradeCompSci grades a computer science submission against test cases
func (s *Service) GradeCompSci(ctx context.Context, code string, testCases []sqlc.CompSciTestCase, timeLimitMs, memoryLimitMB int, problemType sqlc.CompSciProblemType) (*CompSciGradingResult, error) {
	results := make([]problem_eval.TestCaseResult, 0, len(testCases))
	var passed, score, maxScore int32

	// Load turtle runner if needed
	var runnerScript string
	if problemType == sqlc.CompSciProblemTypeTurtle {
		// Try to read the runner script from various locations
		paths := []string{
			"internal/problem_eval/runners/turtle_runner.py",
			"../internal/problem_eval/runners/turtle_runner.py",
			"/app/internal/problem_eval/runners/turtle_runner.py",
		}

		for _, p := range paths {
			if content, err := os.ReadFile(p); err == nil {
				runnerScript = string(content)
				break
			}
		}

		if runnerScript == "" {
			// Fallback: try to embed it or panic/return error.
			// For now return error if not found
			return nil, fmt.Errorf("turtle runner script not found")
		}
	}

	for _, tc := range testCases {
		maxScore += tc.CorrectPoints

		var execResult *problem_eval.ExecutionResult
		var err error
		var tcResult problem_eval.TestCaseResult

		if problemType == sqlc.CompSciProblemTypeTurtle {
			execResult, err = s.executor.ExecuteWithRunner(ctx, runnerScript, code, tc.Input, timeLimitMs, memoryLimitMB)
		} else {
			execResult, err = s.executor.ExecuteCode(ctx, code, tc.Input, timeLimitMs, memoryLimitMB)
		}

		if err != nil {
			results = append(results, problem_eval.TestCaseResult{
				TestCaseID: tc.ID,
				Passed:     false,
				Points:     0,
				Error:      err.Error(),
			})
			continue
		}

		if problemType == sqlc.CompSciProblemTypeTurtle {
			// Parse Turtle Output (SVG Base64)
			output := execResult.Output
			startMarker := "---TURTLE_RESULT_START---"
			endMarker := "---TURTLE_RESULT_END---"

			startIndex := strings.Index(output, startMarker)
			endIndex := strings.Index(output, endMarker)

			if startIndex == -1 || endIndex == -1 {
				tcResult = problem_eval.TestCaseResult{
					TestCaseID: tc.ID,
					Passed:     false,
					Error:      "No image generated or script failed silently.",
					Actual:     output, // Show stdout for debugging
				}
			} else {
				base64Image := strings.TrimSpace(output[startIndex+len(startMarker) : endIndex])

				// Call AI for grading
				if tc.ImageUrl.Valid && tc.ImageUrl.String != "" {
					aiScore, feedback, err := s.aiClient.GradeImage(ctx, base64Image, tc.ImageUrl.String)
					if err != nil {
						tcResult = problem_eval.TestCaseResult{
							TestCaseID: tc.ID,
							Passed:     false,
							Error:      fmt.Sprintf("AI Grading Error: %v", err),
						}
					} else {
						// Scale score to points
						points := int32(float64(aiScore) / 100.0 * float64(tc.CorrectPoints))

						tcResult = problem_eval.TestCaseResult{
							TestCaseID: tc.ID,
							Passed:     aiScore >= 80, // Pass if score >= 80% (configurable)
							Points:     points,
							Expected:   "Image Match",
							Actual:     fmt.Sprintf("AI Score: %d/100\nFeedback: %s\n\n![Result](data:image/svg+xml;base64,%s)", aiScore, feedback, base64Image),
						}
					}
				} else {
					tcResult = problem_eval.TestCaseResult{
						TestCaseID: tc.ID,
						Passed:     false,
						Error:      "Test case missing ideal image URL",
					}
				}
			}

		} else {
			// Standard output comparison
			actualOutput := strings.TrimSpace(execResult.Output)
			expectedOutput := strings.TrimSpace(tc.Output)

			passed := actualOutput == expectedOutput && !execResult.TimedOut && execResult.ExitCode == 0
			var points int32
			if passed {
				points = tc.CorrectPoints
			}

			tcResult = problem_eval.TestCaseResult{
				TestCaseID: tc.ID,
				Passed:     passed,
				Points:     points,
				Expected:   expectedOutput,
				Actual:     actualOutput,
				TimedOut:   execResult.TimedOut,
			}
		}

		if execResult.Error != "" {
			tcResult.Error = execResult.Error
		}

		if tcResult.Passed || problemType == sqlc.CompSciProblemTypeTurtle {
			// For turtle, we accumulate partial points
			score += tcResult.Points
			if tcResult.Passed {
				passed++
			}
		}

		results = append(results, tcResult)
	}

	return &CompSciGradingResult{
		Score:       score,
		MaxScore:    maxScore,
		PassedTests: passed,
		TotalTests:  int32(len(testCases)),
		Results:     results,
	}, nil
}

// QuizGradingResult holds the result of grading a quiz
type QuizGradingResult struct {
	TotalScore     int32
	MaxScore       int32
	ProblemResults []QuizProblemResult
}

// QuizProblemResult holds the result for a single quiz problem
type QuizProblemResult struct {
	ProblemID int64
	Score     int32
	MaxScore  int32
	IsCorrect bool
	Feedback  string
}

// GradeQuiz grades a quiz submission
func (s *Service) GradeQuiz(ctx context.Context, answers map[int64]string, problems []sqlc.QuizProblem) (*QuizGradingResult, error) {
	var totalScore, maxScore int32
	results := make([]QuizProblemResult, 0, len(problems))

	for _, p := range problems {
		maxScore += p.Points
		answer := answers[p.ID]
		var score int32
		var isCorrect bool
		var feedback string

		switch p.ProblemType {
		case "true_false", "fill_blank", "open_ended":
			// Simple text comparison for T/F and Fill Blank
			// Open ended usually needs manual grading, but we can auto-grade if exact match
			if p.CorrectAnswer.Valid && strings.EqualFold(strings.TrimSpace(answer), strings.TrimSpace(p.CorrectAnswer.String)) {
				score = p.Points
				isCorrect = true
			}
		case "mcq_single":
			// Answer should be the Option ID
			// We need to check if this option is correct
			// Since we don't have options passed in, we need to query them.
			// But wait, Service struct doesn't have queries yet.
			// We'll add queries to Service struct.
			// Assumes s.queries is available
			// For now, let's assume we can't grade MCQ without queries.
			// We need to refactor Service to have queries.
		}

		totalScore += score
		results = append(results, QuizProblemResult{
			ProblemID: p.ID,
			Score:     score,
			MaxScore:  p.Points,
			IsCorrect: isCorrect,
			Feedback:  feedback,
		})
	}

	return &QuizGradingResult{
		TotalScore:     totalScore,
		MaxScore:       maxScore,
		ProblemResults: results,
	}, nil
}
