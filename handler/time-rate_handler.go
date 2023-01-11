package handler

import (
	"context"
	"fmt"
	"github.com/BingguWang/limiter-study/global"
	"github.com/BingguWang/limiter-study/limiter"
	"github.com/gin-gonic/gin"
	"time"
)

func TimeRateAllowHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		instance := global.GetMyInstance()

		timeRateLimiter := limiter.NewTimeRateLimiter()
		if !timeRateLimiter.Allow() { // 除了Allow外，还有AllowN,一次性消费N个令牌
			instance.FailedCount++
			fmt.Println("被限流了...:", instance.CurrentTraceIdVal)
			fmt.Println("FailedCount:", instance.FailedCount)
			//ctx.Writer.Write([]byte("被限流了..." + traceId))
			//ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// 每次请求自增一次计数
		instance.SucceedCount++
		fmt.Println("SucceedCount:", instance.SucceedCount)

		// 模拟业务用时

		return
	}
}

// TimeRateWaitHandler wait方法是避免请求返回错误而是给请求等待的时间
func TimeRateWaitHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		instance := global.GetMyInstance()

		c, _ := context.WithTimeout(ctx, time.Millisecond*500) // 最多只会等待这个时间,可以设置时间,请求就会一直等到有令牌为止

		timeRateLimiter := limiter.NewTimeRateLimiter()
		// 没有令牌会在wait()这里阻塞指定的时间，最后wait返回err就说明wait后还是没有令牌可用，就丢弃请求
		if err := timeRateLimiter.Wait(c); err != nil {
			fmt.Println("Error=====: ", err)
			fmt.Println("被限流了...:", instance.CurrentTraceIdVal)
			instance.FailedCount++
			fmt.Println("FailedCount:", instance.FailedCount)
			return
		}

		// 每次请求自增一次计数
		instance.SucceedCount++
		fmt.Println("SucceedCount:", instance.SucceedCount)

		end := time.Now()

		// 模拟业务用时
		sub := end.Sub(start)
		fmt.Println(sub)
		if sub >= 500 {
			fmt.Println("经过wait后succeed: ", instance.CurrentTraceIdVal) // 通过输出的时间可以看到有些请求是阻塞过的,就是wait了的
		}
		return
	}
}

// TimeRateReserveHandler reserve不管token数都返回一个reservation
func TimeRateReserveHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		instance := global.GetMyInstance()
		timeRateLimiter := limiter.NewTimeRateLimiter()
		reservation := timeRateLimiter.Reserve()

		delay := reservation.Delay()
		fmt.Println("需等待的秒数:", delay.Seconds())
		if !reservation.OK() {
			instance.FailedCount++
			fmt.Println("FailedCount:", instance.FailedCount)
			return
		}
		time.Sleep(delay)

		// 每次请求自增一次计数
		instance.SucceedCount++
		fmt.Println("traceId:", instance.CurrentTraceIdVal)
		fmt.Println("SucceedCount:", instance.SucceedCount)

		end := time.Now()
		sub := end.Sub(start)
		fmt.Println(sub.Milliseconds()) // 功能类似wait，推荐使用wait
		return
	}
}
