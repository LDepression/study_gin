package redis

//redis key

//redis key 注意尽量使用命名空间的方式,方便查询和拆分
const (
	KeyPrefix          = "bluebell:"
	KeyPostTimeZSet    = "post:time"  //zset;帖子以发帖时间为分数
	KeyPostScoreZSet   = "post:score" //zset;帖子以投票的分数为分数
	KeyPostVotedZSetPF = "post:voted" //zset记录用户及投票的类型;参数是post id
	KeyCommunitySetPF  = "community:" //set;保存每个分区下帖子的id
)

//给redisKey加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
