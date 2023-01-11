package limiter

import (
	"golang.org/x/time/rate"
	"sync"
)

// -------------- golang time/rate限流器

var (
	myTimeRateLimiterOnce sync.Once
	myTimeRateLimiter     *rate.Limiter
)

func NewTimeRateLimiter() *rate.Limiter {
	myTimeRateLimiterOnce.Do(func() {
		// r 也就是放入令牌的速率，如果传入时间,可以用rate.Every,多久放入一个令牌
		// b 是令牌桶的大小
		myTimeRateLimiter = rate.NewLimiter(500, 200)
	})
	return myTimeRateLimiter
}
