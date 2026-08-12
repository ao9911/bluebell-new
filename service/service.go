package service

import (
	"github.com/ao9911/go-matrix/auth/jwt"

	"github.com/ao9911/bluebell-new/conf"
	"github.com/ao9911/bluebell-new/dao"
)

type Service struct {
	c    *conf.Config
	dao  *dao.Dao
	auth *jwt.JWT
}

func New(c *conf.Config, auth *jwt.JWT) (s *Service) {
	s = &Service{
		c:    c,
		dao:  dao.New(c),
		auth: auth,
	}
	// TODO 初始化其他服务
	return
}

func (s *Service) Close() {

}
