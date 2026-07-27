package dao

import (
	xgorm "github.com/ao9911/go-matrix/database/gorm"
	"gorm.io/gorm"

	"github.com/ao9911/bluebell-new/conf"
)

type Dao struct {
	mysql *gorm.DB
	// TODO 添加其他连接
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		mysql: xgorm.NewORM(c.Mysql),
	}
	return
}
