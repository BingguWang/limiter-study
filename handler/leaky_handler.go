package handler

import (
	"fmt"
	"github.com/BingguWang/limiter-study/global"
	"github.com/BingguWang/limiter-study/utils"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"sync"
	"time"
)

type LeakyBucketLimiter struct {
	rate       float64    // 每秒固定流出速率
	capacity   float64    // 桶的容量
	water      float64    // 当前桶中请求量
	lastLeakMs int64      // 桶上次漏水毫秒数
	lock       sync.Mutex // 锁
}

func (leaky *LeakyBucketLimiter) Allow() bool {
	leaky.lock.Lock()
	defer leaky.lock.Unlock()

	//now := time.Now().UnixNano() / 1e6
	now := time.Now().UnixMilli()
	/**
	假设N是两次流水的时间间隔内流出的请求数
	需保证 rate / 1000 = N / (now - lastLeakMs)
	于是，变形就是 N = (now - lastLeakMs) * rate / 1000
	*/
	// 计算剩余水量,两次执行时间中需要漏掉的水
	leakyWater := leaky.water - (float64(now-leaky.lastLeakMs) * leaky.rate / 1000)
	leaky.water = math.Max(0, leakyWater)
	leaky.lastLeakMs = now
	fmt.Println(now)
	if leaky.water+1 <= leaky.capacity { // 没有满则放入
		leaky.water++
		return true
	} else { // 满了则返回错误
		return false
	}
}

//func (leaky *LeakyBucketLimiter) Set(rate, capacity float64) {
//	leaky.rate = rate
//	leaky.capacity = capacity
//	leaky.water = 0
//	leaky.lastLeakMs = time.Now().UnixNano() / 1e6
//}
func leakyHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		instance := global.GetMyInstance()
		traceId := ctx.Query("traceId")

		// 每秒允许1个请求
		fmt.Println(utils.ToJson(NewLeakyBucketLimiter()))
		if !NewLeakyBucketLimiter().Allow() {
			fmt.Println("被限流了...", traceId)
			fmt.Println("traceId:", traceId)
			instance.FailedCount++
			fmt.Println("SucceedCount:", instance.SucceedCount)
			fmt.Println("FailedCount:", instance.FailedCount)
			ctx.Writer.Write([]byte("被限流了..." + traceId))
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// 每次请求自增一次计数
		instance.SucceedCount++
		fmt.Println("traceId:", traceId)
		fmt.Println("SucceedCount:", instance.SucceedCount)
		fmt.Println("FailedCount:", instance.FailedCount)
		// 模拟业务处理

		return
	}
}
