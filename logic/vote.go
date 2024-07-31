package logic

import (
	"demo_webapp/dao/redis"
	"demo_webapp/models"
	"strconv"
)

func VoteForPost(userID int64, p *models.ParamVote) (err error) {
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
