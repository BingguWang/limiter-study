package handler

import (
	"fmt"
	"github.com/BingguWang/limiter-study/global"
	"github.com/BingguWang/limiter-study/limiter"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func counterHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		instance := global.GetMyInstance()

		// 获取限流器实例
		counterLimiter := limiter.NewCounterLimiter()

		// 判断是否放行该请求
		if !counterLimiter.Allow(instance.CurrentTraceIdVal) {
			//ctx.Writer.Write([]byte("被限流了..." + strconv.Itoa(instance.CurrentTraceIdVal)))
			log.Println("FailedCount:", instance.FailedCount)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// 每次请求自增一次计数
		fmt.Println("SucceedCount:", instance.SucceedCount)

		// 模拟业务处理
		fmt.Println("do your work ...")
		return
	}
}
