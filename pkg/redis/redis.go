package redis

import (
	"context"
	"fmt"
	"gra/pkg/config"

	"github.com/redis/go-redis/v9"
)

func Init(cfg *config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password, // 如果没有密码则留空
		DB:       cfg.DB,       // 默认数据库
		PoolSize: cfg.PoolSize, // 连接池大小
	})

	// 测试连接是否通畅
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("Redis 连接失败: %v", err))
		return nil, err
	}
	return rdb, nil
}
