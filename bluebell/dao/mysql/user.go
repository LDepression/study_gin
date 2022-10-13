package mysql

import (
	"GoAdvance/StudyGinAdvance/bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

//dao层中的函数,待logic层的函数的调用

const secret = "zxz123456.com"

func CheckUserExist(username string) (err error) {
	var id int
	sqlStr := `select id from user where username= ?`
	db.Get(&id, sqlStr, username)
	if id != 0 {
		return ErrorUserExist
	}
	return nil
}

//InsertUser 项数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	//对密码进行加密
	user.Password = encryptPassword(user.Password)
	//执行SQL语句入库
	sqlStr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return nil
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) error {
	oPassword := user.Password
	sqlStr := `select user_id,username,password from user where username = ?`
	err := db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	//查询数据库失败
	if err != nil {
		return err
	}
	//判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return nil
}

func GetUserByID(id int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select
		user_id,username
		from user
		where user_id=?`
	err = db.Get(user, sqlStr, id)
	return
}
