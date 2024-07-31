package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"math"
	"time"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432
)

var (
	ErrVoteTimeExpired = errors.New("超出投票时间")
	ErrVoteRepeated    = errors.New("重复投票")
)

// VoteForPost 投票

func VoteForPost(userID, postID string, value float64) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//1、判断帖子投票时间是否过期
	//先从post:time ZSet中取出post的创建时间

	postTime := client.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpired
	}
	//2、往post:score ZSet中添加分数, 并且更新post:vote:postID ZSet中用户信息，使用pipeline
	//从post:vote:postID中取出当前用户上一次的投票信息（是赞成或反对）

	ov := client.ZScore(ctx, getRedisKey(KeyPostVoteZSetPre+postID), userID).Val()
	if ov == value {
		return ErrVoteRepeated
	}
	//计算需要增加多少分数

	var op float64
	if value-ov > 0 {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(value - ov)

	pipe := client.TxPipeline()
	pipe.ZIncrBy(ctx, getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	if value == 0 {
		pipe.ZRem(ctx, getRedisKey(KeyPostVoteZSetPre+postID), userID)
	} else {
		pipe.ZAdd(ctx, getRedisKey(KeyPostVoteZSetPre+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}

	_, err = pipe.Exec(ctx)
	return
}
