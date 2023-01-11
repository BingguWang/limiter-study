package limiter

import (
	"github.com/BingguWang/limiter-study/global"
	"log"
	"sync"
	"time"
)

// ------------- 令牌桶限流器
/**
以恒定的速率向令牌桶放入令牌,收到请求先取出一个令牌才能处理，无令牌就请求处理错误
【注意的是，我们都无法去控制请求的流入速率，因为请求的发起无法由后端来控制的】
能拿到令牌的请求就能被处理，所以可以处理突发流量，不过如何权衡放令牌和用令牌的速度很重要

令牌桶的容量是固定的,放入令牌速度固定
在请求少的时候，由于恒定的速率放入令牌，就可以攒下来不少的令牌，当有突发流量的时候，令牌可能被取光，取光后就得等恒定速率放令牌，有令牌了就能处理

注意的是
放入令牌的速度不能太慢，这样面对突发流量的时候，不能及时补充令牌回导致大量请求错误

适合电商抢购或者微博出现热点事件这种场景，因为在限流的同时可以应对一定的突发流量。
*/
var (
	myTokenLimiterOnce sync.Once
	myTokenLimiter     *TokenBucketLimiter
)

func NewTokenBucketLimiter() *TokenBucketLimiter {
	myTokenLimiterOnce.Do(func() {
		// 设置桶容量是capacity，放入令牌速率是rate个/s
		myTokenLimiter = &TokenBucketLimiter{rate: 500, capacity: 200}
	})
	return myTokenLimiter
}

type TokenBucketLimiter struct {
	rate         int64 //固定的token放入速率, r/s
	capacity     int64 //桶的容量
	tokens       int64 //桶中当前token数量
	lastTokenSec int64 //上次向桶中放令牌的时间的时间戳，为了精确一点，这里用毫秒
	lock         sync.Mutex
}

func (bucket *TokenBucketLimiter) Take(traceid int) bool {
	bucket.lock.Lock()
	defer bucket.lock.Unlock()
	instance := global.GetMyInstance()

	now := time.Now().UnixMilli()
	/**
	需保证每秒放入令牌数达到rate，pre是放入之前桶内的令牌数, cur是放入后桶内现在有的令牌数
	rate / 1000   = (cur - pre )  / now-bucket.lastTokenSec
	于是
	cur = pre + (now-bucket.lastTokenSec)*rate / 1000
	也就是 bucket.tokens = bucket.tokens + (now-bucket.lastTokenSec)*bucket.rate /1000
	*/
	// 更新桶属性
	bucket.tokens = bucket.tokens + (now-bucket.lastTokenSec)*bucket.rate/1000 // 先添加令牌
	if bucket.tokens >= bucket.capacity {                                      // 桶满，不能再放入令牌
		bucket.tokens = bucket.capacity
	}
	bucket.lastTokenSec = now
	// 桶未空
	if bucket.tokens > 0 {
		bucket.tokens--
		instance.SucceedCount++
		return true
	} else {
		// 桶空,则拒绝
		instance.FailedCount++
		log.Println("被限流了! traceId:  ", traceid)
		return false
	}
}

//func (bucket *TokenBucketLimiter) Init(rate, cap int64) {
//	bucket.rate = rate
//	bucket.capacity = cap
//	bucket.tokens = 0
//	bucket.lastTokenSec = time.Now().Unix()
//}
