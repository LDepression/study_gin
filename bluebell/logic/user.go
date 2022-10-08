package logic

import (
	"GoAdvance/StudyGinAdvance/bluebell/dao/mysql"
	"GoAdvance/StudyGinAdvance/bluebell/models"
	"GoAdvance/StudyGinAdvance/bluebell/pkg/snowflake"
	"fmt"
)

//存放业务逻辑的代码
//之前写书城的时候,我是把没有吧logic这一层单独抽离出来
//我是直接将这一层的logic放到了controller中了,所以显得controller中代码过于乱

func SignUp(p *models.ParamsSignUp) (err error) {
	//1判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		fmt.Println("该用户存在..............")
		return err
	}
	fmt.Println("该用户可以保存")
	//2生成UID
	userID := snowflake.GenID()
	//3.为密码加密
	//构造一个user实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	fmt.Println("user:", user)
	//4.保存进数据库
	if err := mysql.InsertUser(user); err != nil {
		fmt.Println("保存用户失败")
		return err
	}
	return nil

}

//Login 进行登录业务的处理
func Login(user *models.ParamsLogin) error {
	u := &models.User{
		Username: user.Username,
		Password: user.Password,
	}
	return mysql.Login(u)
}
