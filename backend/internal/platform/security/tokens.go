package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// error
var (
	ErrTokenExpired   = errors.New("token has expired")
	ErrTokenInvalid   = errors.New("token is invalid")
	ErrTokenMalformed = errors.New("token is malformed")
)

// payload JWT
type Claims struct {
	UserID    string `json:"uid"`
	Email     string `json:"email"`
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

type TokenConfig struct {
	AccessSecret    string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type TokenProvider interface {
	GenerateAccessToken(userID, email, sessionID string) (string, error)
	ParseAccessToken(tokenString string) (*Claims, error)
	GenerateRefreshToken() (string, error)
	HashToken(token string) string
	RefreshTokenTTL() time.Duration
	AccessTokenTTL() time.Duration
}

type tokenProvider struct {
	cfg TokenConfig
}

func NewTokenProvider(cfg TokenConfig) TokenProvider {
	return &tokenProvider{cfg: cfg}
}

// buat JWT
func (p *tokenProvider) GenerateAccessToken(userID, email, sessionID string) (string, error) {
	now := time.Now()

	claims := &Claims{
		UserID:    userID,
		Email:     email,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(p.cfg.AccessTokenTTL)),
			ID:        uuid.NewString(), // jti — unik per token, berguna untuk blacklist
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.cfg.AccessSecret))
}

// validasi, decode JWT
func (p *tokenProvider) ParseAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrTokenInvalid
			}
			return []byte(p.cfg.AccessSecret), nil
		},
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, ErrTokenMalformed
		default:
			return nil, ErrTokenInvalid
		}
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// token acak 32 byte, di encode ke hex
func (p *tokenProvider) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (p *tokenProvider) HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// expose TTL
func (p *tokenProvider) RefreshTokenTTL() time.Duration {
	return p.cfg.RefreshTokenTTL
}

func (p *tokenProvider) AccessTokenTTL() time.Duration {
	return p.cfg.AccessTokenTTL
}
