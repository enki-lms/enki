package server

import (
	"context"
	"fmt"
	"time"

	"github.com/enki/daemon/internal/ai"
	aiapi "github.com/enki/daemon/internal/api/ai"
	"github.com/enki/daemon/internal/api/courses"
	"github.com/enki/daemon/internal/api/enrollments"
	"github.com/enki/daemon/internal/api/problemgroups"
	"github.com/enki/daemon/internal/api/problems"
	"github.com/enki/daemon/internal/api/submissions"
	"github.com/enki/daemon/internal/api/testcases"
	"github.com/enki/daemon/internal/api/users"
	"github.com/enki/daemon/internal/auth"
	"github.com/enki/daemon/internal/config"
	"github.com/enki/daemon/internal/db"
	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/enki/daemon/internal/problem_eval"
	"github.com/gin-gonic/gin"
)

func StartServer() error {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize database connection
	pool, err := db.NewPoolFromConfig(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	queries := sqlc.New(pool)

	srv := gin.Default()

	// Initialize OIDC provider
	provider, err := auth.NewProvider(ctx, &cfg.OIDC)
	if err != nil {
		return fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.Auth.JWTSecret, 24*time.Hour)

	// Initialize auth middleware
	authMiddleware := auth.NewMiddleware(jwtManager, cfg.Auth.TeacherRoleName, cfg.Auth.AdminEmails, cfg.Sandbox.Unsafe)

	// Initialize sandbox executor (optional - only if nsjail is available)
	var executor *problem_eval.Executor
	executor, err = problem_eval.NewExecutor(&cfg.Sandbox)
	if err != nil {
		fmt.Printf("Warning: Code execution disabled: %v\n", err)
	}

	// Register auth routes
	authHandler := auth.NewHandler(provider, queries, jwtManager, cfg.Server.FrontendURL)
	authGroup := srv.Group("/auth")
	{
		authGroup.GET("/login", authHandler.LoginHandler)
		authGroup.GET("/redirect", authHandler.RedirectHandler)
	}

	// Register API routes
	api := srv.Group("/api")

	// User routes
	userHandler := users.NewHandler(queries, authMiddleware)
	userHandler.RegisterRoutes(api, authMiddleware)

	// Course routes
	courseHandler := courses.NewHandler(queries, authMiddleware)
	courseHandler.RegisterRoutes(api, authMiddleware)

	// Enrollment routes
	enrollmentHandler := enrollments.NewHandler(queries, authMiddleware)
	enrollmentHandler.RegisterRoutes(api, authMiddleware)

	// Problem group routes
	problemGroupHandler := problemgroups.NewHandler(queries, authMiddleware)
	problemGroupHandler.RegisterRoutes(api, authMiddleware)

	// Problem routes
	problemHandler := problems.NewHandler(queries, authMiddleware)
	problemHandler.RegisterRoutes(api, authMiddleware)

	// Test case routes
	testCaseHandler := testcases.NewHandler(queries, authMiddleware)
	testCaseHandler.RegisterRoutes(api, authMiddleware)

	// Submission routes (only if executor is available)
	if executor != nil {
		submissionHandler := submissions.NewHandler(queries, authMiddleware, executor)
		submissionHandler.RegisterRoutes(api, authMiddleware)
	}

	// AI routes (only if enabled)
	aiClient := ai.NewClient(&cfg.AI)
	if aiClient.IsEnabled() {
		aiHandler := aiapi.NewHandler(aiClient)
		aiHandler.RegisterRoutes(api, authMiddleware)
	} else {
		fmt.Println("Info: AI assistant disabled (set ai.enabled=true and ai.api_key in config)")
	}

	return srv.Run(cfg.ServerAddress())
}
