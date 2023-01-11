package handler

import (
	"fmt"
	"github.com/BingguWang/limiter-study/global"
	"github.com/BingguWang/limiter-study/limiter"
	"github.com/gin-gonic/gin"
)

func ChannelHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		instance := global.GetMyInstance()
		channelLimiter := limiter.NewChannelLimiter()
		if !channelLimiter.Allow(instance.CurrentTraceIdVal) {
			//ctx.Writer.Write([]byte("被限流了..." + traceId))
			//ctx.Writer.WriteHeader(http.StatusInternalServerError)
			fmt.Println("FailedCount:", instance.FailedCount)
			return
		}

		fmt.Println("SucceedCount:", instance.SucceedCount)

		// 模拟业务用时
		//time.Sleep(50 * time.Millisecond)

		// 处理完后, 从通道内移除任务
		v, ok := <-channelLimiter.Ch
		if ok {
			fmt.Println(v)
		}
		return
	}
}
