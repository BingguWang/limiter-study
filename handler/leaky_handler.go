package handler

import (
	"fmt"
	"github.com/BingguWang/limiter-study/global"
	"github.com/BingguWang/limiter-study/limiter"
	"github.com/gin-gonic/gin"
)

func LeakyHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		instance := global.GetMyInstance()

		leakyBucketLimiter := limiter.NewLeakyBucketLimiter()

		if !leakyBucketLimiter.Allow(instance.CurrentTraceIdVal) {
			fmt.Println("FailedCount:", instance.FailedCount)
			//ctx.Writer.Write([]byte("被限流了..." + traceId))
			//ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println("SucceedCount:", instance.SucceedCount)

		// 模拟业务处理
		fmt.Println("do your work ...")

		return
	}
}
