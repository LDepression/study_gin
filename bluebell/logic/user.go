package logic

import (
	"GoAdvance/StudyGinAdvance/bluebell/dao/mysql"
	"GoAdvance/StudyGinAdvance/bluebell/pkg/snowflake"
)

//存放业务逻辑的代码
//之前写书城的时候,我是把没有吧logic这一层单独抽离出来
//我是直接将这一层的logic放到了controller中了,所以显得controller中代码过于乱

func SignUp() {
	//1判断用户存不存在
	mysql.QueryUserByUsername()
	//2生成UID
	snowflake.GenID()
	//3.为密码加密
	//4.保存进数据库
	mysql.InsertUser()
	//redis. x
}
