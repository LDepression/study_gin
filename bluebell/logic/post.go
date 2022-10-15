package logic

import (
	"GoAdvance/StudyGinAdvance/bluebell/dao/mysql"
	"GoAdvance/StudyGinAdvance/bluebell/dao/redis"
	"GoAdvance/StudyGinAdvance/bluebell/models"
	"GoAdvance/StudyGinAdvance/bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) error {
	//先生成post_id
	p.ID = snowflake.GenID()
	err := mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID)
	return err
	//再将信息插入数据库中去

}

func GetPostDetailsByID(id int64) (data *models.ApiPostDetails, err error) {
	data = new(models.ApiPostDetails)
	post, err := mysql.GetPostDetailsByID(id)
	data.Post = post
	if err != nil {
		zap.L().Error("mysql.GetPostDetailsByID", zap.Error(err))
		return
	}
	//通过UserID获取作者名
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}
	data.AuthorName = user.Username
	//获取community的详细信息
	communityDetails, err := mysql.GetCommunityDetailList(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityList failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}
	data.CommunityDetail = communityDetails

	return data, err
}

func GetPostList(page, size int64) ([]*models.ApiPostDetails, error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	ApiPostDetails := make([]*models.ApiPostDetails, 0, len(posts))
	for _, post := range posts {
		//通过UserID获取作者名
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//获取community的详细信息
		communityDetails, err := mysql.GetCommunityDetailList(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityList failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetails{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: communityDetails,
		}
		ApiPostDetails = append(ApiPostDetails, postDetail)
	}
	return ApiPostDetails, nil
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetails, err error) {
	//2.去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	//3.根据id去mysql数据库中查询帖子的详细信息
	//返回的数据还要按照给定的顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	//将帖子的作者及分区信息查询出来填充到帖子中去
	for idx, post := range posts {
		//通过UserID获取作者名
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//获取community的详细信息
		communityDetails, err := mysql.GetCommunityDetailList(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityList failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetails{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: communityDetails,
		}
		data = append(data, postDetail)
	}
	return data, nil
}

func GetCommunityPostList(p models.ParamsCommunityPostList) (data []*models.ApiPostDetails, err error) {
	//2.去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	//3.根据id去mysql数据库中查询帖子的详细信息
	//返回的数据还要按照给定的顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	//将帖子的作者及分区信息查询出来填充到帖子中去
	for idx, post := range posts {
		//通过UserID获取作者名
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//获取community的详细信息
		communityDetails, err := mysql.GetCommunityDetailList(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityList failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetails{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: communityDetails,
		}
		data = append(data, postDetail)
	}
	return data, nil
}
