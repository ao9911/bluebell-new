package dao

import (
	"context"
	"errors"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

const refreshTokenKeyPrefix = "auth:refresh:"

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

func (d *Dao) SaveRefreshToken(ctx context.Context, jti string, userID int64, expireSeconds int64) error {
	return d.redisClient.Set(ctx, refreshTokenKey(jti), strconv.FormatInt(userID, 10), time.Duration(expireSeconds)*time.Second).Err()
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

func (d *Dao) RotateRefreshToken(ctx context.Context, oldJTI string, newJTI string, userID int64, expireSeconds int64) error {
	pipeline := d.redisClient.DB().TxPipeline()
	pipeline.Set(ctx, refreshTokenKey(newJTI), strconv.FormatInt(userID, 10), time.Duration(expireSeconds)*time.Second)
	pipeline.Del(ctx, refreshTokenKey(oldJTI))
	_, err := pipeline.Exec(ctx)
	return err
}

func refreshTokenKey(jti string) string {
	return refreshTokenKeyPrefix + jti
}
