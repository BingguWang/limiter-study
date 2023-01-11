package main

import (
	"github.com/BingguWang/limiter-study/global"
	"github.com/BingguWang/limiter-study/limiter"
	"github.com/BingguWang/limiter-study/router"
	"github.com/gin-gonic/gin"
	"log"
	_ "net/http/pprof"
)

var (
	r *gin.Engine
)

func init() {
	// 获取一个全局的单例
	global.GetMyInstance()

	log.Println("开始初始化限流器...")
	limiter.NewCounterLimiter()
	limiter.NewLeakyBucketLimiter()
	limiter.NewTokenBucketLimiter()
	//handler.NewChannelLimiter()
	//handler.NewTimeRateLimiter()

	r = router.SetupRouter()
}

func main() {
	//// 启动一个 http server，注意 pprof 相关的 handler 已经自动注册过了
	//if err := http.ListenAndServe(":6060", nil); err != nil {
	//	log.Fatal(err)
	//}
	if err := r.Run(":8088"); err != nil {
		log.Fatal(err)
	}
}

/**
在实际使用时，一般不会做全局的限流，而是针对某些特征去做精细化的限流。
例如：通过header、x-forward-for 等限制爬虫的访问，通过对 ip,session 等用户信息限制单个用户的访问等。
这里的案例只是简单的全局限流

*/
