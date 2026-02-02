package auth

import (
	"fmt"
	"time"

	"github.com/enki/daemon/internal/db/sqlc/sqlc"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims for authenticated users
type Claims struct {
	jwt.RegisteredClaims
	UserID      int64  `json:"user_id"`
	Email       string `json:"email"`
	Institution string `json:"institution"`
	Role        string `json:"role"`
	FullName    string `json:"full_name"`
}

// JWTManager handles JWT token generation and validation
type JWTManager struct {
	secret        []byte
	tokenDuration time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secret string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secret:        []byte(secret),
		tokenDuration: tokenDuration,
	}
}

// GenerateToken creates a new JWT for the given user
func (m *JWTManager) GenerateToken(user *sqlc.User) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "enki-daemon",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
		UserID:      user.ID,
		Email:       user.Email,
		Institution: user.Institution,
		Role:        string(user.Role),
		FullName:    user.FullName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ValidateToken validates a JWT and returns the claims
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
