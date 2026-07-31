package model

import "time"

type Post struct {
	ID          int64     `gorm:"type:bigint(20);primaryKey;autoIncrement"`
	PostID      int64     `gorm:"type:bigint(20);not null;uniqueIndex:idx_post_id"`
	Title       string    `gorm:"type:varchar(128);not null"`
	Content     string    `gorm:"type:varchar(8192);not null"`
	AuthorID    int64     `gorm:"type:bigint(20);not null;index:idx_author_id"`
	CommunityID int64     `gorm:"type:bigint(20);not null;index:idx_community_id"`
	Status      int8      `gorm:"type:tinyint(4);not null;default:1"`
	CreateTime  time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdateTime  time.Time `gorm:"type:timestamp;autoUpdateTime"`
}

func (Post) TableName() string {
	return "post"
}
