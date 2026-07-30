package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ao9911/bluebell-new/conf"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"

	issuer = "bluebell-new"
	jtiLen = 16
)

var mySecret = []byte("123456")

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidTokenType = errors.New("invalid token type")
)

// MyClaims is the shared JWT claims payload for access and refresh tokens.
type MyClaims struct {
	UserID    int64  `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// GenToken creates a stateless access token and a stateful refresh token.
func GenToken(conf *conf.Config, userID int64) (accessToken, refreshToken, refreshJTI string, err error) {
	accessToken, _, err = genTokenWithJTI(userID, TokenTypeAccess, conf.Auth.AccessExpire)
	if err != nil {
		return "", "", "", err
	}

	refreshToken, refreshJTI, err = genTokenWithJTI(userID, TokenTypeRefresh, conf.Auth.RefreshExpire)
	if err != nil {
		return "", "", "", err
	}

	return accessToken, refreshToken, refreshJTI, nil
}

func genTokenWithJTI(userID int64, tokenType string, expireSeconds int64) (string, string, error) {
	// 生成jti
	buf := make([]byte, jtiLen)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	jti := hex.EncodeToString(buf)

	// 生成token
	claims := MyClaims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireSeconds) * time.Second)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(mySecret)
	if err != nil {
		return "", "", err
	}
	return token, jti, nil
}

// ParseAccessToken parses and validates an access token.
func ParseAccessToken(tokenString string) (*MyClaims, error) {
	return parseTypedToken(tokenString, TokenTypeAccess)
}

// ParseRefreshToken parses and validates a refresh token.
func ParseRefreshToken(tokenString string) (*MyClaims, error) {
	return parseTypedToken(tokenString, TokenTypeRefresh)
}

func parseTypedToken(tokenString, tokenType string) (*MyClaims, error) {
	claims := new(MyClaims)
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}
			return mySecret, nil
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
