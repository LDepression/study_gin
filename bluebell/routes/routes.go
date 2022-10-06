package routes

import (
	concrollers "GoAdvance/StudyGinAdvance/bluebell/controllers"
	"GoAdvance/StudyGinAdvance/bluebell/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	//加入之日库的两个中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//注册业务路由
	r.POST("/signup", concrollers.SignUpHandler)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.Run()
	return r
}
