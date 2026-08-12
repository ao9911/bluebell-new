package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/ao9911/go-matrix/auth/jwt"
	"github.com/ao9911/go-matrix/log"

	v1 "github.com/ao9911/bluebell-new/api/auth/v1"
	"github.com/ao9911/bluebell-new/dao"
	"github.com/ao9911/bluebell-new/model"
	"github.com/ao9911/bluebell-new/pkg/ecode"
	"github.com/ao9911/bluebell-new/pkg/snowflake"
	"github.com/ao9911/bluebell-new/pkg/util"
)

func (s *Service) SignUp(ctx context.Context, param *v1.SignupRequest) (*v1.SignupResponse, error) {
	// 判断用户是否存在
	if err := s.dao.CheckUserExist(ctx, param.Username); err != nil {
		if errors.Is(err, dao.ErrUserExist) {
			return nil, ecode.UserExist
		}
		log.Errorf("s.dao.CheckUserExist error: %v", err)
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
		log.Errorf("s.dao.CreateUser error: %v", err)
		return nil, err
	}
	// 返回结果
	return &v1.SignupResponse{
		UserID:   user.UserID,
		Username: user.Username,
	}, nil
}

func (s *Service) Login(ctx context.Context, auth *jwt.JWT, param *v1.LoginRequest) (*v1.LoginResponse, error) {
	// 查询用户信息
	user, err := s.dao.GetUserByUsername(ctx, param.Username)
	if err != nil {
		if errors.Is(err, dao.ErrUserNotFound) {
			return nil, ecode.UserNotFound
		}
		log.Errorf("s.dao.GetUserByUsername error: %v", err)
		return nil, err
	}
	// 校验密码
	if user.Password != util.EncodeMD5(param.Password) {
		return nil, ecode.InvalidPassword
	}
	// 生成token
	tokenPair, err := auth.GeneratePair(strconv.FormatInt(user.UserID, 10))
	if err != nil {
		log.Errorf("jwt.GenToken error: %v", err)
		return nil, err
	}
	// 解析refresh token
	refreshClaims, err := auth.ParseRefresh(tokenPair.RefreshToken)
	if err != nil {
		log.Errorf("jwt.ParseRefreshToken error: %v", err)
		return nil, err
	}
	// refresh token缓存
	if err := s.dao.SaveRefreshToken(ctx, refreshClaims); err != nil {
		log.Errorf("s.dao.SaveRefreshToken error: %v", err)
		return nil, err
	}
	// 返回结果
	return &v1.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, auth *jwt.JWT, param *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error) {
	// 解析refresh token
	claims, err := auth.ParseRefresh(param.RefreshToken)
	if err != nil {
		log.Errorf("jwt.ParseRefreshToken error: %v", err)
		return nil, ecode.InvalidRefreshToken
	}
	// 校验refresh token是否合法
	userID, err := s.dao.GetRefreshTokenUserID(ctx, claims.ID)
	if err != nil {
		if errors.Is(err, dao.ErrRefreshTokenNotFound) {
			return nil, ecode.InvalidRefreshToken
		}
		log.Errorf("s.dao.GetRefreshTokenUserID error: %v", err)
		return nil, err
	}
	if strconv.FormatInt(userID, 10) != claims.Subject {
		return nil, ecode.InvalidRefreshToken
	}
	// 更新refresh token
	tokenPair, err := auth.GeneratePair(claims.Subject)
	if err != nil {
		log.Errorf("jwt.GenToken error: %v", err)
		return nil, err
	}
	// 解析新的refresh token
	newClaims, err := auth.ParseRefresh(tokenPair.RefreshToken)
	if err != nil {
		log.Errorf("jwt.ParseRefreshToken error: %v", err)
		return nil, err
	}
	// refresh token缓存更新
	if err := s.dao.RotateRefreshToken(ctx, claims, newClaims); err != nil {
		log.Errorf("s.dao.RotateRefreshToken error: %v", err)
		return nil, err
	}
	// 返回结果
	return &v1.RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
