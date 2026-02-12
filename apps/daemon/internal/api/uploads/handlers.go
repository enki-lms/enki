package uploads

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/enki/daemon/internal/auth"
	"github.com/enki/daemon/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles upload-related API routes
type Handler struct {
	cfg *config.ServerConfig
}

// NewHandler creates a new upload handler
func NewHandler(cfg *config.ServerConfig) *Handler {
	// Ensure uploads directory exists
	if err := os.MkdirAll(cfg.UploadsDir, 0755); err != nil {
		fmt.Printf("Warning: Failed to create uploads directory: %v\n", err)
	}
	return &Handler{
		cfg: cfg,
	}
}

// UploadFile uploads a file
// POST /api/uploads
func (h *Handler) UploadFile(c *gin.Context) {
	_, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Single file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	// Validate file type (allow images)
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only image files are allowed"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dst := filepath.Join(h.cfg.UploadsDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Return URL (assume served at /uploads)
	// Build absolute URL if needed, but relative URL /uploads/filename is usually enough for frontend
	// However, LLM might need absolute URL if it can't access local path.
	// But `GradeImage` takes URL. If generic OpenAI API is used, it needs public URL or base64.
	// We are passing URL to `GradeImage`. `GradeImage` sends it to OpenAI.
	// If OpenAI can't access localhost, we have a problem.
	// BUT, `GradeImage` in `client.go` was modifying to accept `studentImageBase64` and `idealImageUrl`.
	// If `idealImageUrl` is a local path (starts with /uploads), OpenAI can't reach it.
	// I should probably pass keys of `idealImage` as base64 too?
	// Or I can read the file in `GradeImage` if it is a local path and convert to base64.
	// Yes, `GradeImage` should handle local paths or the caller should.

	// For now return the relative path.
	url := fmt.Sprintf("/uploads/%s", filename)

	c.JSON(http.StatusOK, gin.H{"url": url})
}

// RegisterRoutes registers upload routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware *auth.Middleware) {
	uploads := rg.Group("/uploads")
	uploads.Use(authMiddleware.AuthRequired())
	{
		uploads.POST("", h.UploadFile)
	}
}
