package v1

import "time"

type CommunityListItem struct {
	CommunityID   uint32 `json:"community_id,string"`
	CommunityName string `json:"community_name"`
}

type CommunityDetail struct {
	CommunityID   uint32    `json:"community_id,string"`
	CommunityName string    `json:"community_name"`
	Introduction  string    `json:"introduction"`
	CreateTime    time.Time `json:"create_time"`
}
