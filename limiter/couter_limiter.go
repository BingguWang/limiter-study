package limiter

import (
	"github.com/BingguWang/limiter-study/global"
	"log"
	"sync"
	"time"
)

// ------------- 计数器限流器
/**
固定时间窗口内对请求进行计数，与阀值进行比较判断是否需要限流，一旦到了时间临界点，将计数器清零
简单说,固定时间段内运行N个请求，但是不关心请求在时间段内的分布情况
限流失败案例：
	比如每一分钟限制100个请求，可以在00:00:00-00:00:58秒里面都没有请求，在00:00:59瞬间发送100个请求，这个对于计数器算法来是允许的，
	然后在00:01:00再次发送100个请求，意味着在短短1s内发送了200个请求
但是由于时间窗每经过1s就会重置计数，就无法识别到这种请求超限，然而这200ms可能导致服务器崩溃
缺点：仍无法应对在时间周期的临界点出现的突发流量，或者说它的痛点就是在时间周期临界点可能出现错误判断而导致限流失败系统崩溃
*/
// 计数器限流器单例
var (
	myCounterLimiterOnce sync.Once
	myCounterLimiter     *CounterLimiter
)

func NewCounterLimiter() *CounterLimiter {
	myCounterLimiterOnce.Do(func() {
		// 每秒允许rate个请求
		myCounterLimiter = &CounterLimiter{rate: 200, cycle: time.Duration(time.Second.Nanoseconds())}
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
	instance := global.GetMyInstance()

	if time.Now().Sub(c.begin) >= c.cycle { // 窗口外
		c.Reset(time.Now())
		instance.SucceedCount++
		return true
	} else { // 窗口内
		if c.rate <= c.count { // 到达阈值
			instance.FailedCount++
			log.Println("被限流了! traceId:  ", traceid)
			return false
		}
		instance.SucceedCount++
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
