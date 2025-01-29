package router

import (
	"github.com/gin-gonic/gin"

	v1 "github.com/AEnjoy/IoT-lubricant/internal/app/core/api/v1"
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/router/middleware"
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

	// v1Route
	v1Route := router.Group("/api/v1")
	signinController := v1.NewAuth()
	v1Route.POST("/signin", signinController.Signin)      // /api/v1/signin
	v1Route.POST("/set-crt", signinController.SetAuthCrt) // /api/v1/set-crt

	routerGroupApp = CommonGroups()
	for _, r := range routerGroupApp {
		r.InitRouter(v1Route, middleware.Auth())
	}

	// Static
	// router.Static("/static", "./static")

	return router, nil
}
