package handler

import (
	"golang.org/x/time/rate"
	"sync"
)

// ------------- 漏桶限流器
var (
	myLeakyLimiterOnce sync.Once
	myLeakyLimiter     *LeakyBucketLimiter
)

func NewLeakyBucketLimiter() *LeakyBucketLimiter {
	myLeakyLimiterOnce.Do(func() {
		// 设置桶容量是capacity，每秒从桶中取出rate个请求处理
		myLeakyLimiter = &LeakyBucketLimiter{rate: 200, capacity: 500}
	})
	return myLeakyLimiter
}

// ------------- 令牌桶限流器

var (
	myTokenLimiterOnce sync.Once
	myTokenLimiter     *TokenBucketLimiter
)

func NewTokenBucketLimiter() *TokenBucketLimiter {
	myTokenLimiterOnce.Do(func() {
		// 设置桶容量是capacity，放入令牌速率是rate个/s
		myTokenLimiter = &TokenBucketLimiter{rate: 200, capacity: 500}
	})
	return myTokenLimiter
}

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
