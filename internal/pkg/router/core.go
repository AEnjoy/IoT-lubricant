package router

import (
	"github.com/AEnjoy/IoT-lubricant/internal/pkg/router/middleware"
	"github.com/gin-gonic/gin"
)

func RouterGroups() []CommonRouter {
	return CommonGroups()
}

var routerGroupApp = RouterGroups()
var middlewares = middleware.GetMiddlewares()

func CoreRouter() (*gin.Engine, error) {
	router := gin.Default()
	router.MaxMultipartMemory = 50 << 20 //50 Gb

	// middleware
	for _, m := range middlewares {
		router.Use(m)
	}

	// Health
	router.GET("/health", Health)

	// v1
	routerGroupApp = CommonGroups()
	privateGroup := router.Group("/api/v1")
	for _, r := range routerGroupApp {
		r.InitRouter(privateGroup)
	}

	// Static
	// router.Static("/static", "./static")

	return router, nil
}
