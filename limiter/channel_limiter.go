package limiter

import (
	"fmt"
	"github.com/BingguWang/limiter-study/global"
	"sync"
)

// -------------- channel限流器
/**
使用缓存通道channel，由于channel是线程安全的，所以 不需要额外的mutex锁
请求来的时候先往channel里塞一个"任务",塞不进说明满了,返回错误
请求处理完，需要取走任务，其实这个原理类似于漏桶，只不过这里没有时间的概念，而且“流出速率”不是恒定的，而是由请求处理完才自己来流出（不过当然是可以通过channel实现漏桶的
*/
var (
	myChannelLimiterOnce sync.Once
	myChannelLimiter     *ChannelLimiter
)

func NewChannelLimiter() *ChannelLimiter {
	myChannelLimiterOnce.Do(func() {
		myChannelLimiter = &ChannelLimiter{Ch: make(chan struct{}, 10)} // 无需设置太大，太大了起不到限流的作用
	})
	return myChannelLimiter
}

type ChannelLimiter struct {
	Ch chan struct{} // 缓存队列用于存放请求
}

func (limit *ChannelLimiter) Allow(traceId int) bool {
	// channel是并发安全的，无需加锁

	instance := global.GetMyInstance()

	select {
	case limit.Ch <- struct{}{}: // 可以放入channel就说明channel没满
		instance.SucceedCount++
		return true
	default:
		fmt.Println("被限流了...", traceId)
		instance.FailedCount++
		return false
	}
}
