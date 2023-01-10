package handler

import (
	"fmt"
	"github.com/BingguWang/limiter-study/global"
	"github.com/BingguWang/limiter-study/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ChannelLimiter struct {
	ch chan struct{} // 缓存队列用于存放请求
}

func (limit *ChannelLimiter) Allow() bool {
	select {
	case limit.ch <- struct{}{}: // 可以放入channel就说明channel没满
		fmt.Println("允许处理")
		return true
	default:
		return false
	}
}

func channelHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		instance := global.GetMyInstance()
		traceId := ctx.Query("traceId")

		fmt.Println(utils.ToJson(NewChannelLimiter()))
		if !NewChannelLimiter().Allow() {
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
		time.Sleep(50 * time.Millisecond)

		// 处理完后
		v, ok := <-NewChannelLimiter().ch
		if ok {
			fmt.Println(v)
		}
		return
	}
}
