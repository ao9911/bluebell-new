package v1

import "time"

// CommunityList.
type CommunityListResponse struct {
	List []*CommunityListItem `json:"list"` // 社区列表
}

type CommunityListItem struct {
	CommunityID   uint32 `json:"community_id,string"`
	CommunityName string `json:"community_name"`
}

// CommunityDetail.
type CommunityDetailResponse struct {
	*CommunityDetail
}

type CommunityDetail struct {
	CommunityID   uint32    `json:"community_id,string"`
	CommunityName string    `json:"community_name"`
	Introduction  string    `json:"introduction"`
	CreateTime    time.Time `json:"create_time"`
}
