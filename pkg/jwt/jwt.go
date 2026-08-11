package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
	jtiLen           = 16
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidTokenType = errors.New("invalid token type")
)

type Config struct {
	AccessSecret  string `toml:"access_secret"`
	RefreshSecret string `toml:"refresh_secret"`
	AccessExpire  int64  `toml:"access_expire"`
	RefreshExpire int64  `toml:"refresh_expire"`
}

// MyClaims is the shared JWT claims payload for access and refresh tokens.
type MyClaims struct {
	UserID    int64  `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// GenToken creates a stateless access token and a stateful refresh token.
func GenToken(c *Config, userID int64) (accessToken, refreshToken, refreshJTI string, err error) {
	accessToken, _, err = genTokenWithJTI(userID, TokenTypeAccess, c.AccessSecret, c.AccessExpire)
	if err != nil {
		return "", "", "", err
	}

	refreshToken, refreshJTI, err = genTokenWithJTI(userID, TokenTypeRefresh, c.RefreshSecret, c.RefreshExpire)
	if err != nil {
		return "", "", "", err
	}

	return accessToken, refreshToken, refreshJTI, nil
}

func genTokenWithJTI(userID int64, tokenType string, secret string, expireSeconds int64) (string, string, error) {
	// 生成jti
	buf := make([]byte, jtiLen)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	jti := hex.EncodeToString(buf)

	// 生成token
	now := time.Now()
	claims := MyClaims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireSeconds) * time.Second)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}
	return token, jti, nil
}

// ParseAccessToken parses and validates an access token.
func ParseAccessToken(c *Config, tokenString string) (*MyClaims, error) {
	return parseTypedToken(tokenString, TokenTypeAccess, c.AccessSecret)
}

// ParseRefreshToken parses and validates a refresh token.
func ParseRefreshToken(c *Config, tokenString string) (*MyClaims, error) {
	return parseTypedToken(tokenString, TokenTypeRefresh, c.RefreshSecret)
}

func parseTypedToken(tokenString, tokenType, secret string) (*MyClaims, error) {
	claims := new(MyClaims)
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	if token == nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.TokenType != tokenType {
		return nil, ErrInvalidTokenType
	}
	return claims, nil
}
