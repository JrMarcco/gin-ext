package rate_limit

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		// develop on wsl, so use the ip of the wsl
		// change it to the ip of the redis server when testing on the your local machine
		Addr: "172.31.176.1:6379",
	})

	return redisClient
}

func TestRedisSlideWindowLimiter(t *testing.T) {
	r := &RedisSlideWindowLimiter{
		Cmd:      initRedis(),
		Interval: time.Second,
		Rate:     1000,
	}

	totalReq := 1500
	var (
		successCnt int
		limitedCnt int
	)

	now := time.Now()
	for range totalReq {
		limited, err := r.Limit(context.Background(), "test_rate_limit_key")
		if err != nil {
			t.Fatalf("failed to limit: %v", err)
			return
		}

		if limited {
			limitedCnt++
		} else {
			successCnt++
		}
	}

	windowStart := now.Add(-time.Second)
	t.Logf("window start at %v", windowStart.Format(time.StampMilli))
	t.Logf("window end at %v", now.Format(time.StampMilli))
	t.Logf("total request count: %d, success count: %d, limited count: %d", totalReq, successCnt, limitedCnt)
}

func TestRedisSlideWindowLimiter_Limit(t *testing.T) {
	r := &RedisSlideWindowLimiter{
		Cmd:      initRedis(),
		Interval: time.Second,
		Rate:     1,
	}

	tcs := []struct {
		name     string
		ctx      context.Context
		key      string
		interval time.Duration
		wantRes  bool
		wantErr  error
	}{
		{
			name:     "success",
			ctx:      context.Background(),
			key:      "test_rate_limit_key",
			interval: time.Second,
			wantRes:  false,
			wantErr:  nil,
		}, {
			name:     "limited",
			ctx:      context.Background(),
			key:      "test_rate_limit_key",
			interval: time.Millisecond * 10,
			wantRes:  true,
			wantErr:  nil,
		}, {
			name:     "window is free",
			ctx:      context.Background(),
			key:      "test_rate_limit_key",
			interval: time.Second,
			wantRes:  false,
			wantErr:  nil,
		}, {
			name:     "another window",
			ctx:      context.Background(),
			key:      "test_another_rate_limit_key",
			interval: time.Millisecond * 10,
			wantRes:  false,
			wantErr:  nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// wait tc.interval time
			<-time.After(tc.interval)
			limited, err := r.Limit(tc.ctx, tc.key)
			assert.Equal(t, tc.wantErr, err)

			if err != nil {
				return
			}

			assert.Equal(t, tc.wantRes, limited)
		})
	}
}
