package v1

import (
	"time"

	communityv1 "github.com/ao9911/bluebell-new/api/community/v1"
)

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// CreatePost.
type CreatePostRequest struct {
	CommunityID int64  `json:"community_id,string" binding:"required"` // 社区id
	Title       string `json:"title" binding:"required"`               // 帖子标题
	Content     string `json:"content" binding:"required"`             // 帖子内容
}

type CreatePostResponse struct {
	PostID int64 `json:"post_id,string"` // 新创建的帖子id
}

// VotePost.
type VotePostRequest struct {
	PostID    int64 `json:"post_id,string" binding:"required"` // 帖子id
	Direction int8  `json:"direction" binding:"oneof=-1 0 1"`  // 投票方向：1赞成，0取消，-1反对
}

// PostList.
type PostListRequest struct {
	CommunityID int64  `json:"community_id" form:"community_id" binding:"omitempty,min=1"`                 // 可以为空，>0 表示按社区过滤
	Page        int64  `json:"page" form:"page,default=1" binding:"min=1" example:"1"`                     // 页码
	Size        int64  `json:"size" form:"size,default=10" binding:"min=1,max=100" example:"10"`           // 每页数据量
	Order       string `json:"order" form:"order,default=time" binding:"oneof=time score" example:"score"` // 排序依据(score或time)
}

type PostListResponse struct {
	List  []*PostListItem `json:"list"`  // 帖子列表
	Total int64           `json:"total"` // 帖子总数
	Page  int64           `json:"page"`  // 当前页码
	Size  int64           `json:"size"`  // 每页数据量
}

type PostListItem struct {
	PostID                       int64              `json:"post_id,string"` // 帖子id
	AuthorName                   string             `json:"author_name"`    // 作者名称
	Title                        string             `json:"title"`          // 帖子标题
	CreateTime                   time.Time          `json:"create_time"`    // 创建时间
	VoteNum                      int64              `json:"vote_num"`       // 投票数
	*communityv1.CommunityDetail `json:"community"` // 社区信息
}

// PostDetail.
type PostDetail struct {
	PostID                       int64              `json:"post_id,string"` // 帖子id
	AuthorName                   string             `json:"author_name"`    // 作者名称
	Title                        string             `json:"title"`          // 帖子标题
	Content                      string             `json:"content"`        // 帖子内容
	CreateTime                   time.Time          `json:"create_time"`    // 创建时间
	*communityv1.CommunityDetail `json:"community"` // 社区信息
}
