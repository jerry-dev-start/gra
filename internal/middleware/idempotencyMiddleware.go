package middleware

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"gra/pkg/response"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func IdempotencyMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// 1. 获取用户 ID (假设已在 Auth 中间件存入 context)
		uid, _ := c.Get("user_id")

		// 2. 读取 Body 内容用于生成哈希
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		// [关键点] 将 Body 重新放回去，供后续 Handler 使用
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 3. 生成 MD5 摘要
		contentHash := fmt.Sprintf("%x", md5.Sum(bodyBytes))

		// 4. 后端拼接唯一的 LockKey
		// 格式：lock:prefix:用户ID:路径:内容哈希
		lockKey := fmt.Sprintf("lock:idmp:%v:%s:%s", uid, c.Request.URL.Path, contentHash)

		// 5. 尝试向 Redis 存入，设置 3 秒过期 (SETNX)
		err := rdb.SetArgs(ctx, lockKey, "1", redis.SetArgs{
			Mode: "NX",
			TTL:  5 * time.Second,
		}).Err()
		if err != nil {
			// 如果 err 等于 redis.Nil，说明 Redis 因为 "NX" 条件（Key已存在）拒绝了写入
			if errors.Is(err, redis.Nil) {
				response.FailMsg(c, "操作过于频繁，请稍后再试")
			} else {
				// 这里的错误是真正的 Redis 连接错误或超时
				response.FailMsg(c, "系统繁忙")
			}
			c.Abort()
			return
		}

		// 如果 err == nil，说明抢锁成功，继续执行
		c.Next()
	}
}
