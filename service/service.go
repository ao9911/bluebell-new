package service

import (
	"github.com/ao9911/bluebell-new/conf"
	"github.com/ao9911/bluebell-new/dao"
)

type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	// TODO 初始化其他服务
	return
}

func (s *Service) Close() {

}
