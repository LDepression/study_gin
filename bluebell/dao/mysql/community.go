package mysql

import (
	"GoAdvance/StudyGinAdvance/bluebell/models"
	"database/sql"

	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id,community_name from community"
	err = db.Select(&communityList, sqlStr)
	if err == sql.ErrNoRows {
		zap.L().Warn("there is no community in db")
		err = nil
	}
	return
}

func GetCommunityDetailList(id int64) (communityDetail *models.CommunityDetail, err error) {
	//这里不能直接用上面的communityDetail,要不然会是空值,一定要去new一个
	communityDetail = new(models.CommunityDetail)
	sqlStr := `select 
			community_id,community_name,introduction,create_time 
			from community
			where community_id=?`
	if err := db.Get(communityDetail, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return communityDetail, err
}
