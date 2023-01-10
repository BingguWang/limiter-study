package handler

import (
	"context"
	"fmt"
	"github.com/BingguWang/limiter-study/global"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Allow方法是直接返回对错结果
func timeRateAllowHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		instance := global.GetMyInstance()
		traceId := ctx.Query("traceId")
		if !NewTimeRateLimiter().Allow() { // 除了Allow外，还有AllowN,一次性消费N个令牌
			fmt.Println("被限流了...", traceId)
			fmt.Println("traceId:", traceId)
			ctx.Writer.Write([]byte("被限流了..." + traceId))
			instance.FailedCount++
			fmt.Println("SucceedCount:", instance.SucceedCount)
			fmt.Println("FailedCount:", instance.FailedCount)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// 每次请求自增一次计数
		instance.SucceedCount++
		fmt.Println("traceId:", traceId)
		fmt.Println("SucceedCount:", instance.SucceedCount)
		fmt.Println("FailedCount:", instance.FailedCount)

		// 模拟业务用时

		return
	}
}

// wait方法是避免请求返回错误而是给请求等待的时间
func timeRateWaitHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		instance := global.GetMyInstance()
		traceId := ctx.Query("traceId")
		c, _ := context.WithTimeout(ctx, time.Millisecond*5) // 最多只会等待这个时间,可以设置时间,请求就会一直等到有令牌为止
		if err := NewTimeRateLimiter().Wait(c); err != nil {
			fmt.Println("Error=====: ", err)
			// 每次请求自增一次计数
			instance.FailedCount++
			fmt.Println("SucceedCount:", instance.SucceedCount)
			fmt.Println("FailedCount:", instance.FailedCount)
			return
		}

		// 每次请求自增一次计数
		instance.SucceedCount++
		fmt.Println("traceId:", traceId)
		fmt.Println("SucceedCount:", instance.SucceedCount)
		fmt.Println("FailedCount:", instance.FailedCount)

		end := time.Now()

		// 模拟业务用时
		sub := end.Sub(start)
		fmt.Println(sub.Milliseconds()) // 通过输出的时间可以看到有些请求是阻塞过的,就是wait了的
		return
	}
}

// reserve不管token数都返回一个reservation
func timeRateReserveHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		instance := global.GetMyInstance()
		traceId := ctx.Query("traceId")
		reservation := NewTimeRateLimiter().Reserve()

		delay := reservation.Delay()
		fmt.Println("需等待的秒数:", delay.Seconds())
		if !reservation.OK() {
			instance.FailedCount++
			fmt.Println("SucceedCount:", instance.SucceedCount)
			fmt.Println("FailedCount:", instance.FailedCount)
			return
		}
		time.Sleep(delay)

		// 每次请求自增一次计数
		instance.SucceedCount++
		fmt.Println("traceId:", traceId)
		fmt.Println("SucceedCount:", instance.SucceedCount)
		fmt.Println("FailedCount:", instance.FailedCount)

		end := time.Now()
		sub := end.Sub(start)
		fmt.Println(sub.Milliseconds()) // 功能类似wait，推荐使用wait
		return
	}
}
