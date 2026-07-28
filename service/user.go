package service

import (
	"context"

	v1 "github.com/ao9911/bluebell-new/api/user/v1"
	"github.com/ao9911/bluebell-new/model"
	"github.com/ao9911/bluebell-new/pkg/ecode"
	"github.com/ao9911/bluebell-new/pkg/jwt"
	"github.com/ao9911/bluebell-new/pkg/snowflake"
	"github.com/ao9911/bluebell-new/pkg/util"
)

func (s *Service) SignUp(ctx context.Context, param *v1.SignupRequest) (*v1.SignupResponse, error) {
	// 判断用户是否存在
	if err := s.dao.CheckUserExist(ctx, param.Username); err != nil {
		return nil, err
	}
	// 创建用户实例
	userID := snowflake.GenID()
	password := util.EncodeMD5(param.Password)
	user := &model.User{
		UserID:   userID,
		Username: param.Username,
		Password: password,
	}
	if err := s.dao.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	// 返回结果
	return &v1.SignupResponse{
		UserID:   user.UserID,
		Username: user.Username,
	}, nil
}

func (s *Service) Login(ctx context.Context, param *v1.LoginRequest) (*v1.LoginResponse, error) {
	// 查询用户信息
	user, err := s.dao.GetUserByUsername(ctx, param.Username)
	if err != nil {
		return nil, err
	}
	// 校验密码
	if user.Password != util.EncodeMD5(param.Password) {
		return nil, ecode.InvalidPassword
	}
	// 生成token
	token, err := jwt.GenToken(s.c, user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	// 返回结果
	return &v1.LoginResponse{
		Token: token,
	}, nil
}
