package ai

import (
	"net/http"

	internalai "github.com/enki/daemon/internal/ai"
	"github.com/enki/daemon/internal/auth"
	"github.com/gin-gonic/gin"
)

// Handler handles AI-related API routes
type Handler struct {
	client *internalai.Client
}

// NewHandler creates a new AI handler
func NewHandler(client *internalai.Client) *Handler {
	return &Handler{
		client: client,
	}
}

// ChatMessageRequest represents a single message in the chat request
type ChatMessageRequest struct {
	Role    string `json:"role" binding:"required,oneof=user assistant"`
	Content string `json:"content" binding:"required"`
}

// ChatRequestBody represents the request body for AI chat
type ChatRequestBody struct {
	Messages []ChatMessageRequest `json:"messages" binding:"required,min=1"`
}

// ChatResponseBody represents the response from AI chat
type ChatResponseBody struct {
	Response string `json:"response"`
}

// Chat sends messages to the AI and returns the response
// POST /api/ai/chat
func (h *Handler) Chat(c *gin.Context) {
	_, ok := auth.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if !h.client.IsEnabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AI assistant is not enabled"})
		return
	}

	var req ChatRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert request messages to AI messages
	messages := make([]internalai.Message, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = internalai.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	response, err := h.client.Chat(c.Request.Context(), messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ChatResponseBody{
		Response: response,
	})
}

// RegisterRoutes registers AI routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware *auth.Middleware) {
	ai := rg.Group("/ai")
	ai.Use(authMiddleware.AuthRequired())
	{
		ai.POST("/chat", h.Chat)
	}
}
