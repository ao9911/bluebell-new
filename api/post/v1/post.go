package v1

import (
	"time"

	communityv1 "github.com/ao9911/bluebell-new/api/community/v1"
)

type CreatePostRequest struct {
	CommunityID int64  `json:"community_id,string" binding:"required"` // 社区id
	Title       string `json:"title" binding:"required"`               // 帖子标题
	Content     string `json:"content" binding:"required"`             // 帖子内容
}

type VotePostRequest struct {
	PostID    int64 `json:"post_id,string" binding:"required"` // 帖子id
	Direction int8  `json:"direction" binding:"oneof=-1 0 1"`  // 投票方向：1赞成，0取消，-1反对
}

type CreatePostResponse struct {
	PostID int64 `json:"post_id,string"` // 新创建的帖子id
}

type PostListItem struct {
	PostID                       int64              `json:"post_id,string"` // 帖子id
	AuthorName                   string             `json:"author_name"`    // 作者名称
	Title                        string             `json:"title"`          // 帖子标题
	CreateTime                   time.Time          `json:"create_time"`    // 创建时间
	*communityv1.CommunityDetail `json:"community"` // 社区信息
}

type PostDetail struct {
	PostID                       int64              `json:"post_id,string"` // 帖子id
	AuthorName                   string             `json:"author_name"`    // 作者名称
	Title                        string             `json:"title"`          // 帖子标题
	Content                      string             `json:"content"`        // 帖子内容
	CreateTime                   time.Time          `json:"create_time"`    // 创建时间
	*communityv1.CommunityDetail `json:"community"` // 社区信息
}
