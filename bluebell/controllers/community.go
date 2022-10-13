package concrollers

import (
	"GoAdvance/StudyGinAdvance/bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//-------社区相关

func CommunityHandle(c *gin.Context) {
	//
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

func CommunityDetailHandle(c *gin.Context) {
	//通过param获取参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	data, err := logic.GetCommunityDetailList(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetailList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
