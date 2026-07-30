package model

import "time"

type Community struct {
	ID            int64     `gorm:"type:int;primaryKey;autoIncrement"`
	CommunityID   uint32    `gorm:"type:int unsigned;not null;uniqueIndex:idx_community_id"`
	CommunityName string    `gorm:"type:varchar(128);not null;uniqueIndex:idx_community_name"`
	Introduction  string    `gorm:"type:varchar(256);not null"`
	CreateTime    time.Time `gorm:"type:timestamp;not null;autoCreateTime"`
	UpdateTime    time.Time `gorm:"type:timestamp;not null;autoUpdateTime"`
}

func (Community) TableName() string {
	return "community"
}
