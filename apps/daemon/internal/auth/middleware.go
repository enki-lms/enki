package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// ContextKeyUser is the key for storing user claims in gin context
	ContextKeyUser = "user"
)

// Middleware holds dependencies for auth middleware
type Middleware struct {
	jwtManager      *JWTManager
	teacherRoleName string
	adminEmails     []string
	unsafe          bool
}

// NewMiddleware creates a new auth middleware
func NewMiddleware(jwtManager *JWTManager, teacherRoleName string, adminEmails []string, unsafe bool) *Middleware {
	return &Middleware{
		jwtManager:      jwtManager,
		teacherRoleName: teacherRoleName,
		adminEmails:     adminEmails,
		unsafe:          unsafe,
	}
}

// AuthRequired validates JWT and sets user claims in context
func (m *Middleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			return
		}

		claims, err := m.jwtManager.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		c.Set(ContextKeyUser, claims)
		c.Next()
	}
}

// RequireRole checks if the user has one of the specified roles
func (m *Middleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get(ContextKeyUser)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "user not authenticated",
			})
			return
		}

		userClaims := claims.(*Claims)

		// Check if user is admin (bypasses role check)
		if m.isAdmin(userClaims.Email) {
			c.Next()
			return
		}

		// Check if user has required role
		for _, role := range roles {
			if userClaims.Role == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "insufficient permissions",
		})
	}
}

// RequireTeacher checks if the user has the teacher role
func (m *Middleware) RequireTeacher() gin.HandlerFunc {
	return m.RequireRole(m.teacherRoleName, "admin")
}

// RequireInstitutionAccess ensures user can only access resources from their institution
// This should be used after AuthRequired() and retrieves institution from route or query
func (m *Middleware) RequireInstitutionAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get(ContextKeyUser)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "user not authenticated",
			})
			return
		}

		userClaims := claims.(*Claims)

		// Admins can access any institution
		if m.isAdmin(userClaims.Email) {
			c.Next()
			return
		}

		// Institution is checked at the handler level since it depends on the resource
		c.Next()
	}
}

// GetUserClaims retrieves user claims from gin context
func GetUserClaims(c *gin.Context) (*Claims, bool) {
	claims, ok := c.Get(ContextKeyUser)
	if !ok {
		return nil, false
	}
	return claims.(*Claims), true
}

// isAdmin checks if the email is in the admin list
func (m *Middleware) isAdmin(email string) bool {
	for _, adminEmail := range m.adminEmails {
		if email == adminEmail {
			return true
		}
	}
	return false
}

// IsTeacherOrAdmin checks if the user has teacher role or is admin
func (m *Middleware) IsTeacherOrAdmin(claims *Claims) bool {
	return claims.Role == m.teacherRoleName || claims.Role == "admin" || m.isAdmin(claims.Email)
}

// CanAccessInstitution checks if user can access resources from the given institution
func (m *Middleware) CanAccessInstitution(claims *Claims, institution string) bool {
	if m.unsafe || m.isAdmin(claims.Email) {
		return true
	}
	return claims.Institution == institution
}
