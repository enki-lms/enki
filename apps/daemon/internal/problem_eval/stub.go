//go:build !linux
// +build !linux

package problem_eval

import (
	"context"
	"fmt"

	"github.com/enki/daemon/internal/config"
)

// ExecutionResult holds the result of running code against a test case
type ExecutionResult struct {
	Output    string `json:"output"`
	Stderr    string `json:"stderr,omitempty"`
	ExitCode  int    `json:"exit_code"`
	TimeTaken int    `json:"time_taken_ms"`
	Error     string `json:"error,omitempty"`
	TimedOut  bool   `json:"timed_out"`
}

// Executor handles sandboxed code execution
type Executor struct {
	cfg *config.SandboxConfig
}

// NewExecutor creates a new executor with the given config
func NewExecutor(cfg *config.SandboxConfig) (*Executor, error) {
	return &Executor{cfg: cfg}, nil
}

// ExecuteCode runs code in a sandbox (stub for non-Linux)
func (e *Executor) ExecuteCode(ctx context.Context, code, input string, timeoutMs, memoryLimitMB int) (*ExecutionResult, error) {
	return nil, fmt.Errorf("code execution is only supported on Linux")
}

// TestCaseResult holds the result of running code against a single test case
type TestCaseResult struct {
	TestCaseID int64  `json:"test_case_id"`
	Passed     bool   `json:"passed"`
	Points     int32  `json:"points"`
	Expected   string `json:"expected,omitempty"`
	Actual     string `json:"actual,omitempty"`
	Error      string `json:"error,omitempty"`
	TimedOut   bool   `json:"timed_out,omitempty"`
}

// SubmissionResult holds the result of running code against all test cases
type SubmissionResult struct {
	ProblemID      int64            `json:"problem_id"`
	TotalTestCases int              `json:"total_test_cases"`
	Passed         int              `json:"passed"`
	Failed         int              `json:"failed"`
	Score          int32            `json:"score"`
	MaxScore       int32            `json:"max_score"`
	Results        []TestCaseResult `json:"results"`
}
