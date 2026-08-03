package dao

import (
	"context"
	"errors"
	"strconv"

	"github.com/ao9911/bluebell-new/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrPostNotFound = errors.New("post not found")

func (d *Dao) CreatePost(ctx context.Context, post *model.Post) error {
	err := d.mysql.WithContext(ctx).Create(post).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao) GetPostByID(ctx context.Context, postID int64) (*model.Post, error) {
	var post model.Post
	err := d.mysql.WithContext(ctx).Where("post_id = ?", postID).First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	return &post, nil
}

func (d *Dao) GetPostList(ctx context.Context, page, size int64) ([]*model.Post, error) {
	var posts []*model.Post
	err := d.mysql.WithContext(ctx).Order("create_time DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (d *Dao) GetPostListByIDs(ctx context.Context, ids []string) ([]*model.Post, error) {
	if len(ids) == 0 {
		return []*model.Post{}, nil
	}

	postIDs := make([]int64, 0, len(ids))
	for _, id := range ids {
		postID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		postIDs = append(postIDs, postID)
	}

	var posts []*model.Post
	err := d.mysql.WithContext(ctx).
		Where("post_id IN ?", postIDs).
		Order(clause.Expr{
			SQL:                "FIELD(post_id,?)",
			Vars:               []interface{}{postIDs},
			WithoutParentheses: true,
		}).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}
