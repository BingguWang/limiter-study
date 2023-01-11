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
		/**
		固定时间窗口内对请求进行计数，与阀值进行比较判断是否需要限流，一旦到了时间临界点，将计数器清零
		简单说,固定时间段内运行N个请求，但是不关心请求在时间段内的分布情况
		限流失败案例：
			比如每一分钟限制100个请求，可以在00:00:00-00:00:58秒里面都没有请求，在00:00:59瞬间发送100个请求，这个对于计数器算法来是允许的，
			然后在00:01:00再次发送100个请求，意味着在短短1s内发送了200个请求
		但是由于时间窗每经过1s就会重置计数，就无法识别到这种请求超限，然而这200ms可能导致服务器崩溃
		缺点：仍无法应对在时间周期的临界点出现的突发流量，或者说它的痛点就是在时间周期临界点可能出现错误判断而导致限流失败系统崩溃
		*/
		group.GET("/counter", handler.CounterHandler())

		// 测试滑动窗口
		/**

		 */
		group.GET("/sliding", handler.SlidingHandler())

		// 测试漏桶
		/**
		最均匀的限流实现方式
		桶的容量是固定的
		流入的请求速率不固定，流出的速率是恒定的,也就是处理请求的速率是恒定的
		桶满的时候(而且桶满的状态是可能会持续一段时间的，所以不是很适合突发流量)，直接返回请求频率超限的错误码或者页面。
		面对突发流量时会有大量请求失败(运行结果可以看到),也就是会有大量的请求被限流而处理失败, 所以不适合电商抢购和微博出现热点事件等场景的限流。
		*/
		group.GET("/leaky", handler.LeakyHandler())

		// 测试令牌桶
		/**
		以恒定的速率向令牌桶放入令牌,收到请求先取出一个令牌才能处理，无令牌就请求处理错误
		【注意的是，我们都无法去控制请求的流入速率，因为请求的发起无法由后端来控制的】
		能拿到令牌的请求就能被处理，所以可以处理突发流量，不过如何权衡放令牌和用令牌的速度很重要

		令牌桶的容量是固定的
		在请求少的时候，由于恒定的速率放入令牌，就可以攒下来不少的令牌，当有突发流量的时候，令牌可能被去光，取光后就得等恒定速率放令牌，有令牌了就能处理
		注意的是放入令牌的速度不能太慢，这样面对突发流量的时候，不能及时补充令牌回导致大量请求错误

		适合电商抢购或者微博出现热点事件这种场景，因为在限流的同时可以应对一定的突发流量。
		*/
		group.GET("/token", handler.TokenHandler())

		// 缓存channel
		/**
		使用缓存通道channel，由于channel是线程安全的，所以 不需要额外的mutex锁
		请求来的时候先往channel里塞一个"门票",塞不进说明满了,返回错误
		请求处理完，需要取走自己的票，其实这个类似于漏桶，只不过这里没有时间的概念，而且“流出速率”不是恒定的，而是由请求处理完自己来流出
		*/
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
		*/

	}

}
