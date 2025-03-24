package rate_limit

import "context"

// Limiter is a rate limiter that can be used to limit the rate of requests.
type Limiter interface {
	// Limit returns true if the request is limited, otherwise returns false.
	Limit(ctx context.Context, key string) (bool, error)
}
