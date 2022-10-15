package routes

import (
	concrollers "GoAdvance/StudyGinAdvance/bluebell/controllers"
	"GoAdvance/StudyGinAdvance/bluebell/logger"
	"GoAdvance/StudyGinAdvance/bluebell/middlewares"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	//加入之日库的两个中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	v1 := r.Group("/api/v1")
	//注册业务路由
	v1.POST("/signup", concrollers.SignUpHandler)
	v1.POST("/login", concrollers.LoginHandler)
	v1.Use(middlewares.JWTAuthMiddleware()) //应用JWT中间件
	{
		v1.GET("/community", concrollers.CommunityHandle)
		v1.GET("/community/:id", concrollers.CommunityDetailHandle)

		v1.POST("/community/post", concrollers.CreatePostHandle)
		v1.GET("/community/post/:id", concrollers.GetPostDetailsByID)
		v1.GET("/posts", concrollers.GetPostListHandle)
		//根据时间或分数获取贴子列表
		v1.GET("/posts2", concrollers.GetPostListHandle2)
		v1.POST("/vote", concrollers.PostVoteController)
	}
	r.Run(":9090")
	return r
}
