package router

import (
	"github.com/BingguWang/limiter-study/global"
	"github.com/BingguWang/limiter-study/handler"
	"github.com/gin-gonic/gin"
)

func TestRouter(r *gin.Engine) {
	r.Use(global.IncrementTraceID)
	group := r.Group("/test")
	{
		// 测试计数器
		group.GET("/counter", handler.CounterHandler())

		// 测试滑动窗口
		/**

		 */
		group.GET("/sliding", handler.SlidingHandler())

		// 测试漏桶
		group.GET("/leaky", handler.LeakyHandler())

		// 测试令牌桶
		group.GET("/token", handler.TokenHandler())

		// 缓存channel
		group.GET("/channel", handler.ChannelHandler())

		// 官方的限流器"golang.org/x/time/rate" [是基于令牌桶实现的]
		// time/rate的Allow的用法,请求没有令牌就直接返回错误
		group.GET("/time-rate-allow", handler.TimeRateAllowHandler())

		// time/rate的wait的用法, 不丢请求，没有令牌时请求可以等待有令牌
		group.GET("/time-rate-wait", handler.TimeRateWaitHandler())

		group.GET("/time-rate-reserve", handler.TimeRateReserveHandler())

		// 分布式限流 Redis + Lua 分布式限流
		/**
		上面的都是单机内的限流，
		分布式限流最关键的是要将限流服务做成原子化
		//TODO 分布式限流
		*/

	}

}
