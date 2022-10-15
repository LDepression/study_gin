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
	page, size := getPageInfo(c)
	data, err := logic.GetPostList(page, size)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

//按创建的 时间 或者 分数进行排序

//GetPostListHandle2 根据前端传来的参数动态获取帖子列表
//1.获取参数
//2.去redis查询id列表
//3.根据id去数据库查询帖子的详细信息
func GetPostListHandle2(c *gin.Context) {
	//get请求参数: /api/v1/posts2?page=1&size=10&order=time  querystring参数
	//先获取分页参数
	//初始化结构体时指定参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime, //magic string
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandle2 failed with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//c.ShouldBind() 动态地获取参数
	//c.ShouldBindJSON() 如果请求参数是json的数据,才能用这个方法获取到
	data, err := logic.GetPostList2(p)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

//GetCommunityPostListHandler 根据社区去查询帖子列表
func GetCommunityPostListHandler(c *gin.Context) {
	//先获取分页参数
	//初始化结构体时指定参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime, //magic string
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandler failed with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//c.ShouldBind() 动态地获取参数
	//c.ShouldBindJSON() 如果请求参数是json的数据,才能用这个方法获取到
	data, err := logic.GetCommunityPostList(p)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
