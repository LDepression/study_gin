package logic

import (
	"GoAdvance/StudyGinAdvance/bluebell/dao/mysql"
	"GoAdvance/StudyGinAdvance/bluebell/models"
	"GoAdvance/StudyGinAdvance/bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) error {
	//先生成post_id
	p.ID = snowflake.GenID()
	//再将信息插入数据库中去
	return mysql.CreatePost(p)
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
