package handler

import (
	"fmt"
	"github.com/BingguWang/limiter-study/global"
	"github.com/BingguWang/limiter-study/limiter"
	"github.com/gin-gonic/gin"
)

func TokenHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		instance := global.GetMyInstance()

		tokenLimiter := limiter.NewTokenBucketLimiter()

		if !tokenLimiter.Take(instance.CurrentTraceIdVal) { // 取令牌
			fmt.Println("FailedCount:", instance.FailedCount)
			//ctx.Writer.Write([]byte("被限流了..." + traceId))
			//ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		instance.SucceedCount++
		fmt.Println("SucceedCount:", instance.SucceedCount)
		// 模拟业务处理

		return
	}
}
