package limiter

import (
	"github.com/BingguWang/limiter-study/global"
	"log"
	"math"
	"sync"
	"time"
)

// ------------- 漏桶限流器
/**
最均匀的限流实现方式

桶的容量是固定的

流入的请求速率不固定
流出的速率是恒定的,也就是处理请求的速率是恒定的

局限性：
桶满的时候(而且桶满的状态是可能会持续一段时间的，所以不是很适合突发流量)，直接返回请求频率超限的错误码或者页面。
面对突发流量时会有大量请求失败(运行结果可以看到),也就是会有大量的请求被限流而处理失败, 所以不适合电商抢购和微博出现热点事件等场景的限流。
*/
var (
	myLeakyLimiterOnce sync.Once
	myLeakyLimiter     *LeakyBucketLimiter
)

func NewLeakyBucketLimiter() *LeakyBucketLimiter {
	myLeakyLimiterOnce.Do(func() {
		// 设置桶容量是capacity，每秒从桶中取出rate个请求处理
		myLeakyLimiter = &LeakyBucketLimiter{rate: 200, capacity: 200}
	})
	return myLeakyLimiter
}

type LeakyBucketLimiter struct {
	rate       float64    // 每秒固定流出速率
	capacity   float64    // 桶的容量
	water      float64    // 当前桶中请求量
	lastLeakMs int64      // 桶上次漏水时间点(毫秒时间戳)
	lock       sync.Mutex // 锁
}

//func (leaky *LeakyBucketLimiter) Set(rate, capacity float64) {
//	leaky.rate = rate
//	leaky.capacity = capacity
//	leaky.water = 0
//	leaky.lastLeakMs = time.Now().UnixNano() / 1e6
//}

func (leaky *LeakyBucketLimiter) Allow(traceid int) bool {
	instance := global.GetMyInstance()
	/**
	假设N是两次流水的时间间隔内流出的请求数
	需保证 rate / 1000 = N / (now - lastLeakMs)
	于是，变形就是 N = (now - lastLeakMs) * rate / 1000
	*/
	// n 是两次是漏水时间之间应该要漏掉多少水
	n := float64(time.Now().UnixMilli()-leaky.lastLeakMs) * leaky.rate / 1000
	// 也就是说当前桶内的水应该有
	rightNowWater := leaky.water - n
	// 更新的leaky.water
	leaky.water = math.Max(0, rightNowWater)
	leaky.lastLeakMs = time.Now().UnixMilli()

	// 判断是否达到阈值
	if leaky.water >= leaky.capacity {
		instance.FailedCount++
		log.Println("被限流了! traceId:  ", traceid)
		return false
	} else {
		instance.SucceedCount++
		leaky.water++
		return true
	}
}
