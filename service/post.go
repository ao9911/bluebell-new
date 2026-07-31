package service

import (
	"context"
	"errors"
	"time"

	"github.com/ao9911/go-matrix/log"

	com_v1 "github.com/ao9911/bluebell-new/api/community/v1"
	v1 "github.com/ao9911/bluebell-new/api/post/v1"
	"github.com/ao9911/bluebell-new/dao"
	"github.com/ao9911/bluebell-new/model"
	"github.com/ao9911/bluebell-new/pkg/ecode"
	"github.com/ao9911/bluebell-new/pkg/snowflake"
)

func (s *Service) CreatePost(ctx context.Context, userID int64, req *v1.CreatePostRequest) (*v1.CreatePostResponse, error) {
	// 检查社区是否存在
	if _, err := s.dao.GetCommunityByID(ctx, req.CommunityID); err != nil {
		if errors.Is(err, dao.ErrCommunityNotFound) {
			return nil, ecode.CommunityNotFound
		}
		log.Errorf("s.dao.GetCommunityByID error: %v", err)
		return nil, err
	}
	// 创建帖子实例
	postID := snowflake.GenID()
	now := time.Now()
	post := &model.Post{
		PostID:      postID,
		Title:       req.Title,
		Content:     req.Content,
		AuthorID:    userID,
		CommunityID: req.CommunityID,
		CreateTime:  now,
	}
	if err := s.dao.CreatePost(ctx, post); err != nil {
		log.Errorf("s.dao.CreatePost error: %v", err)
		return nil, err
	}
	// 帖子ID加入社区set
	if err := s.dao.CreatePostRedisIndex(ctx, post.PostID, post.CommunityID, post.CreateTime); err != nil {
		log.Errorf("s.dao.CreatePostRedisIndex error: %v", err)
		return nil, err
	}
	// 返回结果
	return &v1.CreatePostResponse{
		PostID: postID,
	}, nil
}

func (s *Service) GetPostDetail(ctx context.Context, postID int64) (*v1.PostDetail, error) {
	// 查询帖子详情
	post, err := s.dao.GetPostByID(ctx, postID)
	if err != nil {
		if errors.Is(err, dao.ErrPostNotFound) {
			return nil, ecode.PostNotFound
		}
		log.Errorf("s.dao.GetPostByID error: %v", err)
		return nil, err
	}
	// 查询作者详情
	user, err := s.dao.GetUserByID(ctx, post.AuthorID)
	if err != nil {
		log.Errorf("s.dao.GetUserByID error: %v", err)
		return nil, err
	}
	// 查询社区详情
	community, err := s.dao.GetCommunityByID(ctx, post.CommunityID)
	if err != nil {
		log.Errorf("s.dao.GetCommunityByID error: %v", err)
		return nil, err
	}
	// 返回结果
	return &v1.PostDetail{
		PostID:     post.PostID,
		AuthorName: user.Username,
		Title:      post.Title,
		Content:    post.Content,
		CreateTime: post.CreateTime,
		CommunityDetail: &com_v1.CommunityDetail{
			CommunityID:   community.CommunityID,
			CommunityName: community.CommunityName,
			Introduction:  community.Introduction,
			CreateTime:    community.CreateTime,
		},
	}, nil
}

func (s *Service) GetPostList(ctx context.Context, page, size int64) ([]*v1.PostListItem, error) {
	// 查询帖子列表
	posts, err := s.dao.GetPostList(ctx, page, size)
	if err != nil {
		log.Errorf("s.dao.GetPostList error: %v", err)
		return nil, err
	}
	// 返回结果
	resp := make([]*v1.PostListItem, 0, len(posts))
	for _, p := range posts {
		// 查询作者详情
		user, err := s.dao.GetUserByID(ctx, p.AuthorID)
		if err != nil {
			log.Errorf("s.dao.GetUserByID error: %v", err)
			return nil, err
		}
		// 查询社区详情
		community, err := s.dao.GetCommunityByID(ctx, p.CommunityID)
		if err != nil {
			log.Errorf("s.dao.GetCommunityByID error: %v", err)
			return nil, err
		}
		resp = append(resp, &v1.PostListItem{
			PostID:     p.PostID,
			AuthorName: user.Username,
			Title:      p.Title,
			CreateTime: p.CreateTime,
			CommunityDetail: &com_v1.CommunityDetail{
				CommunityID:   community.CommunityID,
				CommunityName: community.CommunityName,
				Introduction:  community.Introduction,
				CreateTime:    community.CreateTime,
			},
		})
	}
	return resp, nil
}
