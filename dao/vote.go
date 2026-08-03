package dao

import (
	"context"
	"errors"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

const (
	redisKeyPrefix     = "bluebell-new:"
	keyPostTimeZSet    = "post:time"
	keyPostScoreZSet   = "post:score"
	keyCommunitySetPF  = "community:"
	keyPostVotedZSetPF = "post:voted:"
)

func (d *Dao) CreatePostRedisIndex(ctx context.Context, postID, communityID int64, createTime time.Time) error {
	postIDStr := strconv.FormatInt(postID, 10)
	communityIDStr := strconv.FormatInt(communityID, 10)
	score := float64(createTime.Unix())

	pipeline := d.redisClient.DB().TxPipeline()
	// 帖子时间
	pipeline.ZAdd(ctx, getRedisKey(keyPostTimeZSet), goredis.Z{
		Score:  score,
		Member: postIDStr,
	})
	// 帖子分数
	pipeline.ZAdd(ctx, getRedisKey(keyPostScoreZSet), goredis.Z{
		Score:  score,
		Member: postIDStr,
	})
	// 帖子ID加入社区set
	pipeline.SAdd(ctx, getRedisKey(keyCommunitySetPF+communityIDStr), postIDStr)
	_, err := pipeline.Exec(ctx)
	return err
}

func (d *Dao) GetPostIDsInOrder(ctx context.Context, page, size int64, order string) ([]string, error) {
	start, stop := redisPageBounds(page, size)
	return d.redisClient.ZRevRange(ctx, postOrderRedisKey(order), start, stop).Result()
}

func (d *Dao) GetCommunityPostIDsInOrder(ctx context.Context, communityID, page, size int64, order string) ([]string, error) {
	orderKey := postOrderRedisKey(order)
	communityKey := getRedisKey(keyCommunitySetPF + strconv.FormatInt(communityID, 10))
	communityOrderKey := orderKey + ":community:" + strconv.FormatInt(communityID, 10)
	if d.redisClient.DB().Exists(ctx, communityOrderKey).Val() < 1 {
		pipeline := d.redisClient.DB().TxPipeline()
		pipeline.ZInterStore(ctx, communityOrderKey, &goredis.ZStore{
			Keys:      []string{communityKey, orderKey},
			Aggregate: "MAX",
		})
		pipeline.Expire(ctx, communityOrderKey, 60*time.Second)
		if _, err := pipeline.Exec(ctx); err != nil {
			return nil, err
		}
	}

	start, stop := redisPageBounds(page, size)
	return d.redisClient.ZRevRange(ctx, communityOrderKey, start, stop).Result()
}

func (d *Dao) GetPostVoteData(ctx context.Context, ids []string) ([]int64, error) {
	if len(ids) == 0 {
		return []int64{}, nil
	}

	pipeline := d.redisClient.DB().Pipeline()
	commands := make([]*goredis.IntCmd, 0, len(ids))
	for _, id := range ids {
		commands = append(commands, pipeline.ZCount(ctx, getRedisKey(keyPostVotedZSetPF+id), "1", "1"))
	}
	if _, err := pipeline.Exec(ctx); err != nil {
		return nil, err
	}

	votes := make([]int64, 0, len(commands))
	for _, command := range commands {
		votes = append(votes, command.Val())
	}
	return votes, nil
}

func redisPageBounds(page, size int64) (start, stop int64) {
	start = (page - 1) * size
	stop = start + size - 1
	return
}

func postOrderRedisKey(order string) string {
	if order == "score" {
		return getRedisKey(keyPostScoreZSet)
	}
	return getRedisKey(keyPostTimeZSet)
}

func (d *Dao) GetPostVoteDirection(ctx context.Context, userID, postID int64) (float64, error) {
	postIDStr := strconv.FormatInt(postID, 10)
	userIDStr := strconv.FormatInt(userID, 10)

	score, err := d.redisClient.ZScore(ctx, getRedisKey(keyPostVotedZSetPF+postIDStr), userIDStr).Result()
	if errors.Is(err, goredis.Nil) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return score, nil
}

func (d *Dao) UpdatePostVote(ctx context.Context, userID, postID int64, direction int8, scoreDelta float64) error {
	postIDStr := strconv.FormatInt(postID, 10)
	userIDStr := strconv.FormatInt(userID, 10)
	votedKey := getRedisKey(keyPostVotedZSetPF + postIDStr)

	pipeline := d.redisClient.DB().TxPipeline()
	pipeline.ZIncrBy(ctx, getRedisKey(keyPostScoreZSet), scoreDelta, postIDStr)
	if direction == 0 {
		pipeline.ZRem(ctx, votedKey, userIDStr)
	} else {
		pipeline.ZAdd(ctx, votedKey, goredis.Z{
			Score:  float64(direction),
			Member: userIDStr,
		})
	}
	_, err := pipeline.Exec(ctx)
	return err
}

func getRedisKey(key string) string {
	return redisKeyPrefix + key
}
