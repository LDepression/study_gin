package logic

import (
	"GoAdvance/StudyGinAdvance/bluebell/dao/redis"
	"GoAdvance/StudyGinAdvance/bluebell/models"
	"strconv"
)

//投票功能:
/*
1.用户投票的数据
*/

//PostForVote 为帖子投票的函数

//推荐阅读
//基于用户投票的相关算法 https://www.ruanyifeng.com/blog/algorithm/

//本项目使用简化版的投票分数
/*
	投票的几种情况:
direction=1时,有两种情况:
	1.之前没有投反对票                        ->更新分数和投票记录
	2.之前投了反对票,现在改投赞成票             ->更新分数和投票记录
direction=0时,有两种情况                     ->更新分数和投票记录
	1.之前投过赞成票,现在要取消赞成票
	2.之前投过反对票,现在要取消投票
direction=-1时,有两种情况:
	1.之前没有投过票,现在投反对票
	2.之前投赞成票,现在改投反对票
*/

/*
投票的限制:
每个帖子从发表之日起一个星期之内允许用户投票,超过了一个星期就不允许投票了
1.到期之后将redis中保存的赞成票数以及反对的票数储存到mysql中去
2.到期之后删除那个 KeyPostVotedZSetPF
*/

func VoteForPost(userID int64, p *models.ParamsVoteData) error {
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
