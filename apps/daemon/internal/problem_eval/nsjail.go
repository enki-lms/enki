package problem_eval

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/enki/daemon/internal/config"
)

// ExecutionResult holds the result of running code against a test case
type ExecutionResult struct {
	Output    string        `json:"output"`
	Stderr    string        `json:"stderr,omitempty"`
	ExitCode  int           `json:"exit_code"`
	TimeTaken time.Duration `json:"time_taken_ms"`
	Error     string        `json:"error,omitempty"`
	TimedOut  bool          `json:"timed_out"`
}

// Executor handles sandboxed code execution
type Executor struct {
	cfg *config.SandboxConfig
}

// NewExecutor creates a new executor with the given config
func NewExecutor(cfg *config.SandboxConfig) (*Executor, error) {
	// If unsafe mode is enabled, skip environment checks
	if cfg.Unsafe {
		return &Executor{cfg: cfg}, nil
	}

	// Check if nsjail is installed
	if _, err := exec.LookPath(cfg.NSJailPath); err != nil {
		return nil, fmt.Errorf("nsjail not found at %s: %w", cfg.NSJailPath, err)
	}

	// Check if lower directory exists
	if _, err := os.Stat(cfg.LowerDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("lower directory not found at %s: create it using debootstrap", cfg.LowerDir)
	}

	// Ensure temp directory exists
	if err := os.MkdirAll(cfg.TempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	return &Executor{cfg: cfg}, nil
}

// ExecuteCode runs Python code with the given input in a sandbox
func (e *Executor) ExecuteCode(ctx context.Context, code, input string, timeoutMs, memoryLimitMB int) (*ExecutionResult, error) {
	if e.cfg.Unsafe {
		return e.executeUnsafe(ctx, code, input, timeoutMs)
	}
	// Use defaults if not specified
	if timeoutMs <= 0 {
		timeoutMs = e.cfg.TimeoutMs
	}
	if memoryLimitMB <= 0 {
		memoryLimitMB = e.cfg.MemoryLimitMB
	}

	// Create overlay mount for isolation
	overlay, err := NewOverlayMount(e.cfg.TempDir, e.cfg.LowerDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create overlay mount: %w", err)
	}
	defer overlay.Unmount()

	// Mount the overlay
	if err := overlay.Mount(); err != nil {
		return nil, fmt.Errorf("failed to mount overlay: %w", err)
	}

	// Write the user code to a file in the sandbox
	codeFile := filepath.Join(overlay.MergedDir, "code.py")
	if err := os.WriteFile(codeFile, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write code file: %w", err)
	}

	// Build nsjail command
	timeout := time.Duration(timeoutMs) * time.Millisecond
	args := []string{
		"--mode", "o",
		"--time_limit", fmt.Sprintf("%d", (timeoutMs+999)/1000), // Convert to seconds, round up
		"--rlimit_as", fmt.Sprintf("%d", memoryLimitMB),
		"--rlimit_fsize", "10", // 10MB max file size
		"--rlimit_nofile", "32",
		"--rlimit_nproc", "1",
		"--chroot", overlay.MergedDir,
		"--user", "nobody",
		"--group", "nogroup",
		"--disable_proc",
		"--quiet",
		"--", e.cfg.PythonPath, "/code.py",
	}

	// Create context with timeout
	execCtx, cancel := context.WithTimeout(ctx, timeout+500*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(execCtx, e.cfg.NSJailPath, args...)

	// Set up stdin with the test input
	cmd.Stdin = bytes.NewBufferString(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute
	startTime := time.Now()
	err = cmd.Run()
	elapsed := time.Since(startTime)

	result := &ExecutionResult{
		Output:    stdout.String(),
		Stderr:    stderr.String(),
		TimeTaken: elapsed,
		ExitCode:  0,
	}

	if err != nil {
		if execCtx.Err() == context.DeadlineExceeded {
			result.TimedOut = true
			result.Error = "execution timed out"
			result.ExitCode = -1
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.Error = err.Error()
			result.ExitCode = -1
		}
	}

	return result, nil
}

func (e *Executor) executeUnsafe(ctx context.Context, code, input string, timeoutMs int) (*ExecutionResult, error) {
	// Use defaults if not specified
	if timeoutMs <= 0 {
		timeoutMs = e.cfg.TimeoutMs
	}

	// Create a temporary file for the code
	tmpFile, err := os.CreateTemp("", "code-*.py")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(code); err != nil {
		return nil, fmt.Errorf("failed to write code to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temp file: %w", err)
	}

	// Create context with timeout
	timeout := time.Duration(timeoutMs) * time.Millisecond
	execCtx, cancel := context.WithTimeout(ctx, timeout+500*time.Millisecond)
	defer cancel()

	// Run python directly
	// Note: We use "python3" assuming it's in the PATH. Safe mode uses cfg.PythonCmd but that might be a path inside the chroot.
	// For unsafe mode, we just want to run it on the host.
	cmd := exec.CommandContext(execCtx, "python3", tmpFile.Name())

	// Set up stdin
	cmd.Stdin = bytes.NewBufferString(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	startTime := time.Now()
	err = cmd.Run()
	elapsed := time.Since(startTime)

	result := &ExecutionResult{
		Output:    stdout.String(),
		Stderr:    stderr.String(),
		TimeTaken: elapsed,
		ExitCode:  0,
	}

	if err != nil {
		if execCtx.Err() == context.DeadlineExceeded {
			result.TimedOut = true
			result.Error = "execution timed out"
			result.ExitCode = -1
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.Error = err.Error()
			result.ExitCode = -1
		}
	}

	return result, nil
}

// TestCaseResult holds the result of a single test case execution
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
