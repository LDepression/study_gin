package concrollers

import (
	"GoAdvance/StudyGinAdvance/bluebell/logic"
	"GoAdvance/StudyGinAdvance/bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreatePostHandle(c *gin.Context) {
	p := new(models.Post)
	//1.获取参数即参数的校验
	//2.shouldbindjson //validator -->binding tag
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("c.ShouldBindJSON(p) failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//这里将作者的id传给p
	AuthorID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = AuthorID

	//2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		//不要把服务端的错误传给用户端
		ResponseError(c, CodeServerBusy)
		return

	}
	ResponseSuccess(c, nil)
}

func GetPostDetailsByID(c *gin.Context) {
	//先获取id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	data, err := logic.GetPostDetailsByID(id)
	if err != nil {
		zap.L().Error("logic.GetPostDetailsByID failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

//GetPostListHandle 获取帖子列表参数
func GetPostListHandle(c *gin.Context) {
	//先获取分页参数
	var (
		size int64
		page int64
		err  error
	)
	size, err = strconv.ParseInt(c.Query("size"), 10, 64)
	if err != nil {
		size = 10
	}
	page, err = strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil {
		page = 1
	}

	data, err := logic.GetPostList(page, size)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
