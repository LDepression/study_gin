package mysql

import (
	"GoAdvance/StudyGinAdvance/web_app/settings"
	"fmt"

	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init(cfg *settings.MySQLConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName)
	db, err := sqlx.Connect("mysql", dsn)

	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		//fmt.Printf("connext db failed,err:%v\n", err)
		return err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return nil
}

func Close() {
	db.Close()
}
