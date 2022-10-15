package redis

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis/v8"
)

//本项目使用简化版的投票分数
/*
	投票的几种情况:
direction=1时,有两种情况:
	1.之前没有投反对票                        ->更新分数和投票记录 差值的绝对值:1 +432
	2.之前投了反对票,现在改投赞成票             ->更新分数和投票记录 差值的绝对值:2 +432*2
direction=0时,有两种情况
	1.之前投过赞成票,现在要取消赞成票								差值的绝对值:1 +432
	2.之前投过反对票,现在要取消投票								差值的绝对值:1 -432
direction=-1时,有两种情况:
	1.之前没有投过票,现在投反对票								差值的绝对值:1 -432
	2.之前投赞成票,现在改投反对票								差值的绝对值:2 -432*2
*/
const (
	oneWeekInSeconds = 7 * 24 * 1600
	scorePerVote     = 432 //每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
)

func CreatePost(postID int64) error {
	ctx := context.Background()
	//帖子的时间
	pipe := client.TxPipeline()
	pipe.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//帖子的分数
	pipe.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	_, err := pipe.Exec(ctx)
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	ctx := context.Background()
	//1.判断投票的限制
	//去redis去帖子发布的时间
	postTime := client.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	//2和3需要放到一个pipeline事务中去

	//2. 更新帖子的分数
	//先查当前的用户帖子投票记录
	ov := client.ZScore(ctx, getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	var dir float64
	if value > ov {
		dir = 1
	} else {
		dir = -1
	}
	diff := math.Abs(ov - value) //计算两次投票的差值
	pipe := client.TxPipeline()
	pipe.ZIncrBy(ctx, getRedisKey(KeyPostScoreZSet), dir*diff*scorePerVote, postID)
	//3. 记录用户为该帖子投票的数据

	if value == 0 {
		pipe.ZRem(ctx, getRedisKey(KeyPostVotedZSetPF+postID), postID)
	} else {
		pipe.ZAdd(ctx, getRedisKey(KeyPostVotedZSetPF+postID), &redis.Z{
			Score:  value, //赞成票还是反对票
			Member: userID,
		})
	}
	_, err := pipe.Exec(ctx)
	return err
}
