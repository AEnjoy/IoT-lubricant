package router

import (
	"os"

	def "github.com/aenjoy/iot-lubricant/pkg/constant"
	v1 "github.com/aenjoy/iot-lubricant/services/apiserver/api"
	"github.com/aenjoy/iot-lubricant/services/apiserver/router/middleware"
	"github.com/gin-gonic/gin"
)

var middlewares = middleware.GetMiddlewares()

func CoreRouter() (*gin.Engine, error) {
	if os.Getenv(def.ENV_RUNNING_LEVEL) != "debug" {
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
	v1Route.POST("/refresh-token", signinController.RefreshToken)
	v1Route.GET("/signin", signinController.Signin)       // /api/v1/signin (oauth callback)
	v1Route.POST("/login", signinController.Login)        // /api/v1/login
	v1Route.POST("/set-crt", signinController.SetAuthCrt) // /api/v1/set-crt

	v1Route.Use(middleware.Auth())
	routerGroupApp := CommonGroups
	for _, r := range routerGroupApp {
		r.InitRouter(v1Route)
	}

	v1Route.GET("/get-private-key", signinController.GetPrivateKey)
	// Static
	// router.Static("/static", "./static")

	return router, nil
}
