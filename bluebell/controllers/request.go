package concrollers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

//getCurrentUser 获取当前登录的用户ID
func getCurrentUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

func getPageInfo(c *gin.Context) (int64, int64) {
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
	return page, size
}
