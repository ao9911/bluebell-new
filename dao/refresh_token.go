package dao

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/ao9911/go-matrix/auth/jwt"
	goredis "github.com/redis/go-redis/v9"
)

const refreshTokenKeyPrefix = "auth:refresh:"

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

func (d *Dao) SaveRefreshToken(ctx context.Context, claims *jwt.Claims) error {
	return d.redisClient.Set(ctx, refreshTokenKey(claims.ID), claims.Subject, time.Until(claims.ExpiresAt.Time)).Err()
}

func (d *Dao) GetRefreshTokenUserID(ctx context.Context, jti string) (int64, error) {
	userID, err := d.redisClient.Get(ctx, refreshTokenKey(jti)).Result()
	if errors.Is(err, goredis.Nil) {
		return 0, ErrRefreshTokenNotFound
	}
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(userID, 10, 64)
}

func (d *Dao) DeleteRefreshToken(ctx context.Context, jti string) error {
	return d.redisClient.Del(ctx, refreshTokenKey(jti)).Err()
}

func (d *Dao) RotateRefreshToken(ctx context.Context, oldClaims, newClaims *jwt.Claims) error {
	pipeline := d.redisClient.DB().TxPipeline()
	pipeline.Set(ctx, refreshTokenKey(newClaims.ID), newClaims.Subject, time.Until(newClaims.ExpiresAt.Time))
	pipeline.Del(ctx, refreshTokenKey(oldClaims.ID))
	_, err := pipeline.Exec(ctx)
	return err
}

func refreshTokenKey(jti string) string {
	return refreshTokenKeyPrefix + jti
}
