package exam

import (
	"context"
	"log"
	"time"

	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/enki/daemon/internal/websocket"
	"github.com/jackc/pgx/v5/pgtype"
)

// TimerService checks for expired exam sessions and auto-submits
type TimerService struct {
	queries  *sqlc.Queries
	hub      *websocket.Hub
	interval time.Duration
	stopCh   chan struct{}
}

// NewTimerService creates a new timer service
func NewTimerService(queries *sqlc.Queries, hub *websocket.Hub) *TimerService {
	return &TimerService{
		queries:  queries,
		hub:      hub,
		interval: 5 * time.Second, // Check every 5 seconds
		stopCh:   make(chan struct{}),
	}
}

// Start starts the timer service
func (s *TimerService) Start() {
	go s.run()
	log.Println("Exam timer service started")
}

// Stop stops the timer service
func (s *TimerService) Stop() {
	close(s.stopCh)
	log.Println("Exam timer service stopped")
}

func (s *TimerService) run() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkExpiredSessions()
		case <-s.stopCh:
			return
		}
	}
}

func (s *TimerService) checkExpiredSessions() {
	ctx := context.Background()

	// Find students whose time has expired but haven't submitted
	students, err := s.queries.ListExamStudentsNeedingAutoSubmit(ctx)
	if err != nil {
		log.Printf("Failed to list students needing auto-submit: %v", err)
		return
	}

	for _, student := range students {
		s.autoSubmit(ctx, student)
	}
}

func (s *TimerService) autoSubmit(ctx context.Context, student sqlc.ExamSessionStudent) {
	log.Printf("Auto-submitting exam for student ID %d (session student ID %d)", student.UserID, student.ID)

	// Get their work in progress
	work, err := s.queries.ListExamWorkInProgressByStudent(ctx, student.ID)
	if err != nil {
		log.Printf("Failed to get work in progress for student %d: %v", student.ID, err)
	}

	// Get the session to know the problem type
	session, err := s.queries.GetExamSession(ctx, student.SessionID)
	if err != nil {
		log.Printf("Failed to get session for student %d: %v", student.ID, err)
		return
	}

	// Create submissions for each problem's work in progress
	for _, w := range work {
		if session.ProblemGroupType == "comp_sci" {
			// For comp_sci problems, we need to evaluate the code
			// For now, just save as-is with 0 score (needs actual evaluation)
			_, err = s.queries.CreateCompSciSubmission(ctx, sqlc.CreateCompSciSubmissionParams{
				UserID:      student.UserID,
				ProblemID:   w.ProblemID,
				Code:        w.CurrentAnswer,
				Score:       0,
				MaxScore:    0, // Will be calculated by evaluator
				PassedTests: 0,
				TotalTests:  0,
				ResultsJson: `{"auto_submitted": true}`,
			})
			if err != nil {
				log.Printf("Failed to create comp_sci submission for student %d, problem %d: %v",
					student.ID, w.ProblemID, err)
			}
		} else if session.ProblemGroupType == "quiz" {
			// For quiz problems, create a submission (answer is JSON)
			_, err = s.queries.CreateQuizSubmission(ctx, sqlc.CreateQuizSubmissionParams{
				UserID:     student.UserID,
				ProblemID:  w.ProblemID,
				AnswerText: pgtype.Text{String: w.CurrentAnswer, Valid: w.CurrentAnswer != ""},
				Score:      0,
				MaxScore:   0, // Will need to be set from problem
			})
			if err != nil {
				log.Printf("Failed to create quiz submission for student %d, problem %d: %v",
					student.ID, w.ProblemID, err)
			}
		}
	}

	// Mark as submitted
	_, err = s.queries.SubmitExamSession(ctx, sqlc.SubmitExamSessionParams{
		ID:            student.ID,
		AutoSubmitted: pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		log.Printf("Failed to mark student %d as submitted: %v", student.ID, err)
		return
	}

	// Notify the client if still connected
	s.hub.EndSessionForStudent(student.SessionID, student.UserID, "time_expired")

	// Clean up work in progress
	err = s.queries.DeleteExamWorkInProgressByStudent(ctx, student.ID)
	if err != nil {
		log.Printf("Failed to delete work in progress for student %d: %v", student.ID, err)
	}

	log.Printf("Auto-submit complete for student ID %d", student.ID)
}
