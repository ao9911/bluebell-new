package dao

import (
	"context"
	"errors"

	"github.com/ao9911/bluebell-new/model"
	"gorm.io/gorm"
)

var ErrCommunityNotFound = errors.New("community not found")

func (d *Dao) GetCommunityList(ctx context.Context) ([]*model.Community, error) {
	// 查询社区列表
	var communities []*model.Community
	err := d.mysql.WithContext(ctx).Find(&communities).Error
	if err != nil {
		return nil, err
	}
	return communities, nil
}

func (d *Dao) GetCommunityByID(ctx context.Context, communityID int64) (*model.Community, error) {
	// 查询社区详情
	var community model.Community
	err := d.mysql.WithContext(ctx).Where("community_id = ?", communityID).First(&community).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommunityNotFound
		}
		return nil, err
	}
	return &community, nil
}
