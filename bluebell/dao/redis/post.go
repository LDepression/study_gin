package redis

import (
	"GoAdvance/StudyGinAdvance/bluebell/models"
	"context"

	"github.com/go-redis/redis/v8"
)

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis中获取id
	//根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	//确定查询的索引的起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	ctx := context.Background()
	//3.zRevRange查询,按分数从大到小,查询指定数量的分数
	return client.ZRevRange(ctx, key, start, end).Result()
}

//GetPostVoteData 根据id查询没票帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {

	//下面这种做法没有问题,但是效率不高
	//如果数据大的话,每一次都要发一次请求
	//我们可以用pipeline,将多次查询按一次请求发过去

	//data := make([]int64, 0, len(ids))
	//ctx := context.Background()
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPF + id)
	//	//查找key中分数是1的元素的数量 -->统计每篇帖子的赞成票的数量
	//	v1 := client.ZCount(ctx, key, "1", "1").Val()
	//	data = append(data, v1)
	//}

	//使用pipeline一次发送多条命令,减少RTT
	ctx := context.Background()
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

//GetCommunityPostIDsInOrder 按社区查询ids
func GetCommunityPostIDsInOrder() ([]string, error) {
	//从redis中获取id
	//根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	//确定查询的索引的起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	ctx := context.Background()
	//3.zRevRange查询,按分数从大到小,查询指定数量的分数
	return client.ZRevRange(ctx, key, start, end).Result()
}
