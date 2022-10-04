package redis

import (
	"GoAdvance/StudyGinAdvance/web_app/settings"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func Init(cfg *settings.RedisConfig) error {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port),
		Password: cfg.Password, // 密码
		DB:       cfg.DB,       // 数据库
		PoolSize: cfg.PoolSize, // 连接池大小
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	return err
}
func Close() {
	_ = rdb.Close()
}
