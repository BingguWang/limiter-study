package global

import (
	"github.com/gin-gonic/gin"
)

func IncrementTraceID(ctx *gin.Context) {
	// 考虑加锁
	GetMyInstance().CurrentTraceIdVal++
}
