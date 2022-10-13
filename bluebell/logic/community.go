package logic

import (
	"GoAdvance/StudyGinAdvance/bluebell/dao/mysql"
	"GoAdvance/StudyGinAdvance/bluebell/models"
)

func GetCommunityList() ([]*models.Community, error) {
	//查数据库,查找到所有的community,然后返回
	return mysql.GetCommunityList()
}

func GetCommunityDetailList(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailList(id)
}
