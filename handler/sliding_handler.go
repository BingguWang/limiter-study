package handler

import (
	"fmt"
	"github.com/BingguWang/limiter-study/global"
	"github.com/gin-gonic/gin"
)

func SlidingHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		instance := global.GetMyInstance()
		traceId := ctx.Query("traceId")
		// 每次请求自增一次计数
		instance.SucceedCount++
		fmt.Println("traceId:", traceId)
		fmt.Println("SucceedCount:", instance.SucceedCount)
	}
}
