package handler

import (
	"golang.org/x/time/rate"
	"sync"
)

// -------------- channel限流器

var (
	myChannelLimiterOnce sync.Once
	myChannelLimiter     *ChannelLimiter
)

func NewChannelLimiter() *ChannelLimiter {
	myChannelLimiterOnce.Do(func() {
		myChannelLimiter = &ChannelLimiter{ch: make(chan struct{}, 10)} // 无需设置太大，太大了起不到限流的作用
	})
	return myChannelLimiter
}

// -------------- golang time/rate限流器

var (
	myTimeRateLimiterOnce sync.Once
	myTimeRateLimiter     *rate.Limiter
)

func NewTimeRateLimiter() *rate.Limiter {
	myTimeRateLimiterOnce.Do(func() {
		// r 也就是放入令牌的速率，如果传入时间,可以用rate.Every,多久放入一个令牌
		// b 是令牌桶的大小
		myTimeRateLimiter = rate.NewLimiter(50, 200)
	})
	return myTimeRateLimiter
}
