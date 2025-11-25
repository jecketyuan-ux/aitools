package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/eduflow/eduflow/internal/pkg/response"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	rdb      *redis.Client
	duration int
	limit    int
}

func NewRateLimiter(rdb *redis.Client, duration, limit int) *RateLimiter {
	return &RateLimiter{
		rdb:      rdb,
		duration: duration,
		limit:    limit,
	}
}

func (r *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		key := fmt.Sprintf("rate_limit:%s", c.ClientIP())

		count, err := r.rdb.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			r.rdb.Expire(ctx, key, time.Duration(r.duration)*time.Second)
		}

		if count > int64(r.limit) {
			response.RateLimitExceeded(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
