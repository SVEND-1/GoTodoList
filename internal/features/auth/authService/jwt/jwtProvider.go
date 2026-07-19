package jwt

import (
	"TodoList/internal/core/domain"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
)

type Claims struct {
	UserID int         `json:"user_id"`
	Role   domain.Role `json:"role"`
	Type   string      `json:"type"`
	jwt.RegisteredClaims
}

type JwtProvider struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJwtProvider(secret string, accessTTL, refreshTTL time.Duration) *JwtProvider {
	return &JwtProvider{
		secret:          []byte(secret),
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

func (p *JwtProvider) GenerateAccessToken(userID int, role domain.Role) (string, error) {
	return p.generateToken(userID, role, "access", p.accessTokenTTL)
}

func (p *JwtProvider) GenerateRefreshToken(userID int, role domain.Role) (string, error) {
	return p.generateToken(userID, role, "refresh", p.refreshTokenTTL)
}

func (s *JwtProvider) generateToken(userID int, role domain.Role, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)

	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

func (s *JwtProvider) ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *JwtProvider) ParseAccessToken(tokenString string) (*Claims, error) {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.Type != "access" {
		return nil, fmt.Errorf("%w: expected access token", ErrInvalidToken)
	}
	return claims, nil
}

func (s *JwtProvider) ParseRefreshToken(tokenString string) (*Claims, error) {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.Type != "refresh" {
		return nil, fmt.Errorf("%w: expected refresh token", ErrInvalidToken)
	}
	return claims, nil
}

func (p *JwtProvider) GetTokenFromCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
