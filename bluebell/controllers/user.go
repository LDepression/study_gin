package concrollers

import (
	"GoAdvance/StudyGinAdvance/bluebell/logic"
	"GoAdvance/StudyGinAdvance/bluebell/models"
	"net/http"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

//SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	//1.获取参数和参数校验
	p := new(models.ParamsSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errs.Translate(trans)), //翻译错误
		})
		return
	}
	//手动对传入的参数进行校验
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 {
	//	zap.L().Error("SignUp with invalid param")
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": "请求参数有误",
	//	})
	//}
	//2.业务处理
	logic.SignUp()
	//3.返回响应
	c.JSON(http.StatusOK, gin.H{
		"status": "传入参数正确",
	})
}
