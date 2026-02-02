package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

const (
	stateCookieName = "oidc_state"
	stateCookieAge  = 300 // 5 minutes
)

// Handler handles OIDC authentication routes
type Handler struct {
	provider    *Provider
	queries     *sqlc.Queries
	jwtManager  *JWTManager
	frontendURL string
}

// NewHandler creates a new auth handler
func NewHandler(provider *Provider, queries *sqlc.Queries, jwtManager *JWTManager, frontendURL string) *Handler {
	return &Handler{
		provider:    provider,
		queries:     queries,
		jwtManager:  jwtManager,
		frontendURL: frontendURL,
	}
}

// generateState creates a cryptographically secure random state token
func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// LoginHandler initiates the OIDC login flow
// GET /auth/login
func (h *Handler) LoginHandler(c *gin.Context) {
	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate state",
		})
		return
	}

	// Store state in a secure cookie for CSRF protection
	c.SetCookie(
		stateCookieName,
		state,
		stateCookieAge,
		"/",
		"",    // domain - empty for current domain
		false, // secure - set to true in production with HTTPS
		true,  // httpOnly
	)

	// Redirect to OIDC provider
	authURL := h.provider.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// RedirectHandler handles the OIDC callback
// GET /auth/redirect
func (h *Handler) RedirectHandler(c *gin.Context) {
	// Check for error response from provider
	if errParam := c.Query("error"); errParam != "" {
		errDesc := c.Query("error_description")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       errParam,
			"description": errDesc,
		})
		return
	}

	// Validate state parameter
	state := c.Query("state")
	storedState, err := c.Cookie(stateCookieName)
	if err != nil || state == "" || state != storedState {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid state parameter",
		})
		return
	}

	// Clear the state cookie
	c.SetCookie(stateCookieName, "", -1, "/", "", false, true)

	// Exchange authorization code for tokens
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing authorization code",
		})
		return
	}

	token, err := h.provider.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to exchange code for tokens",
		})
		return
	}

	userInfo, err := h.provider.UserInfo(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get user info",
		})
		return
	}

	// Parse userInfo claims
	var userClaims map[string]interface{}
	if err = userInfo.Claims(&userClaims); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to parse user claims",
		})
		return
	}
	fmt.Println(userClaims)

	// Extract user info from claims
	email, ok := userClaims["email"].(string)
	if !ok || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email not found in user claims",
		})
		return
	}

	// Extract displayName from list (first element)
	displayName := extractFirstFromList(userClaims, "displayName")
	if displayName == "" {
		displayName = email // Fallback to email if displayName not available
	}

	// Extract hrEduPersonAffiliation from list (first element)
	affiliation := extractFirstFromList(userClaims, "hrEduPersonAffiliation")

	// Check if user exists
	ctx := c.Request.Context()
	user, err := h.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// User doesn't exist, create them
			user, err = h.queries.CreateUser(ctx, sqlc.CreateUserParams{
				Email:       email,
				Institution: affiliation,
				FullName:    displayName,
				GivenName:   extractFirstFromList(userClaims, "givenName"),
				Role:        sqlc.UserRoleStudent, // Default role
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "failed to create user",
					"details": err.Error(),
				})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to query user",
				"details": err.Error(),
			})
			return
		}
	}

	// Generate JWT token for the user
	jwtToken, err := h.jwtManager.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	// Set token cookie
	c.SetCookie(
		"token",
		jwtToken,
		3600*24, // 24 hours
		"/",
		"",
		false, // secure (true in prod)
		false, // httpOnly (false so JS can read it if needed, or true if using strict cookie auth)
		// For this demo, let's keep it false or rely on the header.
		// Actually, to make it easy for the frontend to pick up without an API call, we can set it not HttpOnly?
		// Better practice: HttpOnly, but we need a way to know we are logged in.
		// Let's set it as HttpOnly = false so the frontend can read it and put it in local storage or use it.
		// Or better, set it HttpOnly=false for this simple demo.
	)

	c.Redirect(http.StatusTemporaryRedirect, h.frontendURL)
}

// extractFirstFromList extracts the first string element from a claim that is a list
func extractFirstFromList(claims map[string]interface{}, key string) string {
	value, ok := claims[key]
	if !ok {
		return ""
	}

	// If it's already a string, return it
	if str, ok := value.(string); ok {
		return str
	}

	// If it's a list, get the first element
	if list, ok := value.([]interface{}); ok && len(list) > 0 {
		if str, ok := list[0].(string); ok {
			return str
		}
	}

	return ""
}
