package router

import (
	"github.com/BingguWang/limiter-study/handler"
	"github.com/gin-gonic/gin"
)

type RouteWork func(*gin.Engine)

var routeWorkSlice = []RouteWork{}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	injectRouteWork(handler.TestRouter)
	for _, work := range routeWorkSlice {
		work(r)
	}
	return r
}

// 注入路由工作
func injectRouteWork(rwork ...RouteWork) {
	routeWorkSlice = append(routeWorkSlice, rwork...)
}
