package dao

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/ao9911/bluebell-new/model"
)

var (
	ErrUserExist    = errors.New("user exist")
	ErrUserNotFound = errors.New("user not found")
)

// 判断用户是否存在
func (d *Dao) CheckUserExist(ctx context.Context, username string) error {
	var count int64
	err := d.mysql.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrUserExist
	}
	return nil
}

// 创建用户
func (d *Dao) CreateUser(ctx context.Context, user *model.User) error {
	err := d.mysql.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

// 查询用户信息（名字）
func (d *Dao) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := d.mysql.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		// 判断是否存在记录
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// 查询用户信息（ID）
func (d *Dao) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	var user model.User
	err := d.mysql.WithContext(ctx).Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		// 判断是否存在记录
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
