package rate_limit

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

//go:embed sliding_window.lua
var redisSlidingWindowLua string

// RedisSlideWindowLimiter is a limiter that uses a Redis sorted set to implement a sliding window rate limit.
// Interval is the size of the window, and Rate is the threshold of the rate limit.
// Means that in the interval all the requests are limited to the Rate.
type RedisSlidingWindowLimiter struct {
	Cmd      redis.Cmdable
	Interval time.Duration // window sieze
	Rate     int64         // threshold
}

func NewRedisSlidingWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int64) *RedisSlidingWindowLimiter {
	return &RedisSlidingWindowLimiter{
		Cmd:      cmd,
		Interval: interval,
		Rate:     rate,
	}
}

// Limit checks if the request is limited.
// If the request is limited, it returns true.
// If the request is not limited, it returns false.
func (r *RedisSlidingWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return false, fmt.Errorf("failed to generate uuid: %w", err)
	}

	// if the result is 0, the request is limited
	result, err := r.Cmd.Eval(
		ctx,
		redisSlidingWindowLua,
		[]string{key},
		r.Interval.Milliseconds(),
		r.Rate,
		time.Now().UnixMilli(),
		uid.String(),
	).Int()

	if err != nil {
		return false, fmt.Errorf("failed to evaluate redis lua script at redis sliding window limiter: %w", err)
	}

	return result == 0, nil
}
