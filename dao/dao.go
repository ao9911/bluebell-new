package dao

import (
	xredis "github.com/ao9911/go-matrix/cache/redis"
	xgorm "github.com/ao9911/go-matrix/database/gorm"
	"gorm.io/gorm"

	"github.com/ao9911/bluebell-new/conf"
)

type Dao struct {
	mysql       *gorm.DB
	redisClient *xredis.RedisStorage
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		mysql:       xgorm.NewORM(c.Mysql),
		redisClient: xredis.NewRedisClient(c.Redis),
	}
	return
}
