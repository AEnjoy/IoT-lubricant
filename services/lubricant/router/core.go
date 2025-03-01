package router

import (
	"os"

	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1"
	"github.com/aenjoy/iot-lubricant/services/lubricant/router/middleware"
	"github.com/gin-gonic/gin"
)

var middlewares = middleware.GetMiddlewares()

func CoreRouter() (*gin.Engine, error) {
	if os.Getenv("RUNNING_LEVEL") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
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
	v1Route.GET("/signin", signinController.Signin)       // /api/v1/signin (oauth callback)
	v1Route.POST("/login", signinController.Login)        // /api/v1/login
	v1Route.POST("/set-crt", signinController.SetAuthCrt) // /api/v1/set-crt

	v1Route.Use(middleware.Auth())
	routerGroupApp := CommonGroups
	for _, r := range routerGroupApp {
		r.InitRouter(v1Route)
	}

	// Static
	// router.Static("/static", "./static")

	return router, nil
}
