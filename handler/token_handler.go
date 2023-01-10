package handler

import (
	"fmt"
	"github.com/BingguWang/limiter-study/global"
	"github.com/BingguWang/limiter-study/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

type TokenBucketLimiter struct {
	rate         int64 //固定的token放入速率, r/s
	capacity     int64 //桶的容量
	tokens       int64 //桶中当前token数量
	lastTokenSec int64 //上次向桶中放令牌的时间的时间戳，单位为秒
	lock         sync.Mutex
}

func (bucket *TokenBucketLimiter) Take() bool {
	bucket.lock.Lock()
	defer bucket.lock.Unlock()

	now := time.Now().Unix()
	/**
	需保证每秒放入令牌数达到rate，pre是放入之前桶内的令牌数, cur是放入后桶内现在有的令牌数
	rate   = (cur - pre )  / now-bucket.lastTokenSec
	cur = pre + (now-bucket.lastTokenSec)
	也就是 bucket.tokens = bucket.tokens + (now-bucket.lastTokenSec)*bucket.rate
	*/
	bucket.tokens = bucket.tokens + (now-bucket.lastTokenSec)*bucket.rate // 先添加令牌
	if bucket.tokens > bucket.capacity {                                  // 桶满，不能再放入令牌
		bucket.tokens = bucket.capacity
	}
	bucket.lastTokenSec = now
	if bucket.tokens > 0 {
		// 还有令牌，领取令牌
		bucket.tokens--
		return true
	} else {
		// 没有令牌,则拒绝
		return false
	}
}

//func (bucket *TokenBucketLimiter) Init(rate, cap int64) {
//	bucket.rate = rate
//	bucket.capacity = cap
//	bucket.tokens = 0
//	bucket.lastTokenSec = time.Now().Unix()
//}
func tokenHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		instance := global.GetMyInstance()
		traceId := ctx.Query("traceId")

		// 每秒允许1个请求
		fmt.Println(utils.ToJson(NewTokenBucketLimiter()))
		if !NewTokenBucketLimiter().Take() { // 取令牌
			fmt.Println("被限流了...", traceId)
			fmt.Println("traceId:", traceId)
			instance.FailedCount++
			fmt.Println("SucceedCount:", instance.SucceedCount)
			fmt.Println("FailedCount:", instance.FailedCount)
			ctx.Writer.Write([]byte("被限流了..." + traceId))
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// 每次请求自增一次计数
		instance.SucceedCount++
		fmt.Println("traceId:", traceId)
		fmt.Println("SucceedCount:", instance.SucceedCount)
		fmt.Println("FailedCount:", instance.FailedCount)
		// 模拟业务处理

		return
	}
}
