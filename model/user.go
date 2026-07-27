package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID   int64  `gorm:"type:bigint(20);not null;uniqueIndex:idx_user_id"`
	Username string `gorm:"type:varchar(64);not null;uniqueIndex:idx_username"`
	Password string `gorm:"type:varchar(64);not null"`
}

func (User) TableName() string {
	return "user"
}
