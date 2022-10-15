package concrollers

import (
	"GoAdvance/StudyGinAdvance/bluebell/logic"
	"GoAdvance/StudyGinAdvance/bluebell/models"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

//投票

func PostVoteController(c *gin.Context) {
	//参数校验
	p := new(models.ParamsVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) //类型断言
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) //翻译并去除掉错误提示中的结构体标签
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
	}
	//获取当前请求的用户id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	//具体投票业务逻辑
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost(userID, p) failed",
			zap.Int64("userID", userID),
			zap.String("postID", p.PostID),
			zap.Int8("direction", int8(p.Direction)))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
