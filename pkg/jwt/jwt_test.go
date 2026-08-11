package jwt

import "testing"

var config *Config

func init() {
	config = &Config{
		AccessSecret:  "access_secret",
		RefreshSecret: "refresh_secret",
		AccessExpire:  3600,
		RefreshExpire: 7200,
	}
}

func TestGenToken(t *testing.T) {
	accessToken, refreshToken, refreshJTI, err := GenToken(config, 1)
	if err != nil {
		t.Fatalf("GenToken() error = %v", err)
	}
	if accessToken == "" {
		t.Fatal("GenToken() access token is empty")
	}
	if refreshToken == "" {
		t.Fatal("GenToken() refresh token is empty")
	}
	if refreshJTI == "" {
		t.Fatal("GenToken() refresh JTI is empty")
	}
}

func TestParseToken(t *testing.T) {
	accessToken, refreshToken, _, err := GenToken(config, 1)
	if err != nil {
		t.Fatalf("GenToken() error = %v", err)
	}

	accessClaims, err := ParseAccessToken(config, accessToken)
	if err != nil {
		t.Fatalf("ParseAccessToken() error = %v", err)
	}
	if accessClaims.UserID != 1 || accessClaims.TokenType != TokenTypeAccess {
		t.Fatalf("access claims = %+v, want user_id 1 and type %q", accessClaims, TokenTypeAccess)
	}

	refreshClaims, err := ParseRefreshToken(config, refreshToken)
	if err != nil {
		t.Fatalf("ParseRefreshToken() error = %v", err)
	}
	if refreshClaims.UserID != 1 || refreshClaims.TokenType != TokenTypeRefresh {
		t.Fatalf("refresh claims = %+v, want user_id 1 and type %q", refreshClaims, TokenTypeRefresh)
	}
}
