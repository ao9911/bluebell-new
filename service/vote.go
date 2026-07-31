package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/ao9911/go-matrix/log"

	v1 "github.com/ao9911/bluebell-new/api/post/v1"
	"github.com/ao9911/bluebell-new/dao"
	"github.com/ao9911/bluebell-new/pkg/ecode"
)

const (
	voteExpireSeconds = 7 * 24 * 3600
	scorePerVote      = 432
)

func (s *Service) PostVote(ctx context.Context, userID int64, req *v1.VotePostRequest) error {
	// 查询帖子详情
	post, err := s.dao.GetPostByID(ctx, req.PostID)
	if err != nil {
		if errors.Is(err, dao.ErrPostNotFound) {
			return ecode.PostNotFound
		}
		log.Errorf("s.dao.GetPostByID error: %v", err)
		return err
	}
	// 检查投票时间是否过期
	if time.Since(post.CreateTime) > time.Duration(voteExpireSeconds)*time.Second {
		return ecode.VoteTimeExpired
	}
	// 查询用户对该帖子的投票方向
	oldDirection, err := s.dao.GetPostVoteDirection(ctx, userID, req.PostID)
	if err != nil {
		log.Errorf("s.dao.GetPostVoteDirection error: %v", err)
		return err
	}
	if float64(req.Direction) == oldDirection {
		return ecode.VoteRepeated
	}
	// 计算分数增量
	diff := math.Abs(oldDirection - float64(req.Direction))
	op := -1.0
	if float64(req.Direction) > oldDirection {
		op = 1.0
	}
	scoreDelta := op * diff * scorePerVote
	// 更新帖子分数和用户投票记录
	if err := s.dao.UpdatePostVote(ctx, userID, req.PostID, req.Direction, scoreDelta); err != nil {
		log.Errorf("s.dao.UpdatePostVote error: %v", err)
		return err
	}
	return nil
}
