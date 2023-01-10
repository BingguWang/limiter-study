package handler

import (
	"github.com/BingguWang/limiter-study/global"
	"github.com/gin-gonic/gin"
)

func IncrementTraceID(ctx *gin.Context) {
	// 考虑加锁
	global.GetMyInstance().CurrentTraceIdVal++
}
