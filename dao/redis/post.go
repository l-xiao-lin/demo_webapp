package redis

import (
	"context"
	"demo_webapp/models"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func CreatePost(postID, communityID int64) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//1、往post:time添加数据
	pipe := client.TxPipeline()
	pipe.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	//2、往post:score添加数据

	pipe.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	//3、往community:id set中添加数据
	cKey := getRedisKey(KeyCommunitySetPre + strconv.Itoa(int(communityID)))
	pipe.SAdd(ctx, cKey, postID)

	_, err = pipe.Exec(ctx)
	return

}

// GetPostIDsInOrder 获取帖子IDS
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	return client.ZRevRange(ctx, key, start, end).Result()
}

// GetCommunityPostIDsInOrder 基于社区获取帖子IDS
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//1、将community set 与order ZSet进行interSort 生成新的一ZSet
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	cKey := getRedisKey(KeyCommunitySetPre + strconv.Itoa(int(p.CommunityID)))
	key := orderKey + strconv.Itoa(int(p.CommunityID))

	if client.Exists(ctx, key).Val() < 1 {
		pipe := client.TxPipeline()
		pipe.ZInterStore(ctx, key, &redis.ZStore{
			Keys:      []string{orderKey, cKey},
			Aggregate: "MAX",
		})
		pipe.Expire(ctx, key, 60*time.Minute)
		_, err := pipe.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	//2、从新的ZSet中获取帖子IDS
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	return client.ZRevRange(ctx, key, start, end).Result()

}

// GetPostVoteData 获取帖子的投票数
func GetPostVoteData(ids []string) (data []int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipe := client.Pipeline()
	for _, id := range ids {
		pipe.ZCount(ctx, getRedisKey(KeyPostVoteZSetPre+id), "1", "1")
	}
	cmders, err := pipe.Exec(ctx)
	if err != nil {
		return
	}

	//断言取值
	data = make([]int64, 0, len(ids))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}
