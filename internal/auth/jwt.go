package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type Claims struct {
	Subject      int64
	Email        string
	IsSuperadmin bool
	TokenType    string
	IssuedAt     time.Time
	ExpiresAt    time.Time
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type JWTManager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	now        func() time.Time
}

func NewJWTManager(secret string, accessTTL, refreshTTL time.Duration) *JWTManager {
	return &JWTManager{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		now:        time.Now,
	}
}

func (m *JWTManager) IssuePair(userID int64, email string, isSuperadmin bool) (*TokenPair, error) {
	now := m.now().UTC()
	access, err := m.signToken(Claims{
		Subject:      userID,
		Email:        email,
		IsSuperadmin: isSuperadmin,
		TokenType:    "access",
		IssuedAt:     now,
		ExpiresAt:    now.Add(m.accessTTL),
	})
	if err != nil {
		return nil, err
	}
	refresh, err := m.signToken(Claims{
		Subject:      userID,
		Email:        email,
		IsSuperadmin: isSuperadmin,
		TokenType:    "refresh",
		IssuedAt:     now,
		ExpiresAt:    now.Add(m.refreshTTL),
	})
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "Bearer",
		ExpiresIn:    int64(m.accessTTL.Seconds()),
	}, nil
}

func (m *JWTManager) VerifyAccess(token string) (*Claims, error) {
	c, err := m.verify(token)
	if err != nil {
		return nil, err
	}
	if c.TokenType != "access" {
		return nil, ErrInvalidToken
	}
	return c, nil
}

func (m *JWTManager) VerifyRefresh(token string) (*Claims, error) {
	c, err := m.verify(token)
	if err != nil {
		return nil, err
	}
	if c.TokenType != "refresh" {
		return nil, ErrInvalidToken
	}
	return c, nil
}

func (m *JWTManager) signToken(c Claims) (string, error) {
	header := map[string]any{"alg": "HS256", "typ": "JWT"}
	payload := map[string]any{
		"sub":           c.Subject,
		"email":         c.Email,
		"is_superadmin": c.IsSuperadmin,
		"token_type":    c.TokenType,
		"iat":           c.IssuedAt.Unix(),
		"exp":           c.ExpiresAt.Unix(),
	}
	h, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	p, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	unsigned := encodeBase64URL(h) + "." + encodeBase64URL(p)
	sig := m.sign(unsigned)
	return unsigned + "." + encodeBase64URL(sig), nil
}

func (m *JWTManager) verify(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}
	unsigned := parts[0] + "." + parts[1]
	sigRaw, err := decodeBase64URL(parts[2])
	if err != nil {
		return nil, ErrInvalidToken
	}
	expected := m.sign(unsigned)
	if !hmac.Equal(sigRaw, expected) {
		return nil, ErrInvalidToken
	}
	payloadBytes, err := decodeBase64URL(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}
	payload := map[string]any{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, ErrInvalidToken
	}

	claims := &Claims{}
	if claims.Subject, err = asInt64(payload["sub"]); err != nil {
		return nil, ErrInvalidToken
	}
	if claims.Email, err = asString(payload["email"]); err != nil {
		return nil, ErrInvalidToken
	}
	if claims.IsSuperadmin, err = asBool(payload["is_superadmin"]); err != nil {
		return nil, ErrInvalidToken
	}
	if claims.TokenType, err = asString(payload["token_type"]); err != nil {
		return nil, ErrInvalidToken
	}
	iat, err := asInt64(payload["iat"])
	if err != nil {
		return nil, ErrInvalidToken
	}
	exp, err := asInt64(payload["exp"])
	if err != nil {
		return nil, ErrInvalidToken
	}
	claims.IssuedAt = time.Unix(iat, 0).UTC()
	claims.ExpiresAt = time.Unix(exp, 0).UTC()
	if m.now().UTC().After(claims.ExpiresAt) {
		return nil, ErrExpiredToken
	}
	return claims, nil
}

func (m *JWTManager) sign(input string) []byte {
	h := hmac.New(sha256.New, m.secret)
	_, _ = h.Write([]byte(input))
	return h.Sum(nil)
}

func encodeBase64URL(v []byte) string {
	return base64.RawURLEncoding.EncodeToString(v)
}

func decodeBase64URL(v string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(v)
}

func asInt64(v any) (int64, error) {
	switch x := v.(type) {
	case float64:
		return int64(x), nil
	case int64:
		return x, nil
	case int:
		return int64(x), nil
	case json.Number:
		return x.Int64()
	default:
		return 0, fmt.Errorf("not an int64")
	}
}

func asString(v any) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("not a string")
	}
	return s, nil
}

func asBool(v any) (bool, error) {
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("not a bool")
	}
	return b, nil
}
