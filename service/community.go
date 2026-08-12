package service

import (
	"context"
	"errors"

	v1 "github.com/ao9911/bluebell-new/api/community/v1"
	"github.com/ao9911/bluebell-new/dao"
	"github.com/ao9911/bluebell-new/pkg/ecode"
	"github.com/ao9911/go-matrix/log"
)

func (s *Service) GetCommunityList(ctx context.Context) (*v1.CommunityList, error) {
	// 查询社区列表
	communities, err := s.dao.GetCommunityList(ctx)
	if err != nil {
		log.Errorf("s.dao.GetCommunityList error: %v", err)
		return nil, err
	}
	// 返回结果
	resp := make([]*v1.CommunityListItem, 0, len(communities))
	for _, c := range communities {
		resp = append(resp, &v1.CommunityListItem{
			CommunityID:   c.CommunityID,
			CommunityName: c.CommunityName,
		})
	}
	return &v1.CommunityList{
		List: resp,
	}, nil
}

func (s *Service) GetCommunityDetail(ctx context.Context, communityID int64) (*v1.CommunityDetail, error) {
	// 查询社区详情
	community, err := s.dao.GetCommunityByID(ctx, communityID)
	if err != nil {
		if errors.Is(err, dao.ErrCommunityNotFound) {
			return nil, ecode.CommunityNotFound
		}
		log.Errorf("s.dao.GetCommunityByID error: %v", err)
		return nil, err
	}
	// 返回结果
	resp := &v1.CommunityDetail{
		CommunityID:   community.CommunityID,
		CommunityName: community.CommunityName,
		Introduction:  community.Introduction,
		CreateTime:    community.CreateTime,
	}
	return resp, nil

}
