package redis

const (
	Prefix             = "demo_webapp:"
	KeyPostTimeZSet    = "post:time"
	KeyPostScoreZSet   = "post:score"
	KeyPostVoteZSetPre = "post:vote:"
	KeyCommunitySetPre = "community:"
)

func getRedisKey(key string) string {
	return Prefix + key
}
