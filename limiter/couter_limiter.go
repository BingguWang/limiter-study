package limiter

import (
	"log"
	"sync"
	"time"
)

// ------------- 计数器限流器
// 计数器限流器单例
var (
	myCounterLimiterOnce sync.Once
	myCounterLimiter     *CounterLimiter
)

func NewCounterLimiter() *CounterLimiter {
	myCounterLimiterOnce.Do(func() {
		// 每秒允许rate个请求
		myCounterLimiter = &CounterLimiter{rate: 500, cycle: time.Duration(time.Second.Nanoseconds())}
	})
	return myCounterLimiter
}

type CounterLimiter struct {
	rate  int           // 阀值,允许的最大请求数
	begin time.Time     // 计数开始时间
	cycle time.Duration // 计数周期
	count int           // 收到的请求数
	lock  sync.Mutex    // 锁
}

// Allow 拦截逻辑，是否放行
func (c *CounterLimiter) Allow(traceid int) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	// 先判断是否达到阈值
	if c.rate < c.count {
		// 判断是否在窗口内
		if time.Now().Sub(c.begin) >= c.cycle {
			c.Reset(time.Now())
			return true
		} else {
			log.Println("被限流了! traceId:  ", traceid)
			return false
		}
	} else {
		c.count++
		return true
	}
}


//func (limit *CounterLimiter) Set(rate int, cycle time.Duration) {
//	limit.rate = rate
//	limit.begin = time.Now()
//	limit.cycle = cycle
//	limit.count = 0
//}

func (limit *CounterLimiter) Reset(begin time.Time) {
	limit.begin = begin
	limit.count = 0
}
