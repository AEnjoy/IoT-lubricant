package middleware

import "github.com/gin-gonic/gin"

func GetMiddlewares() []func(ctx *gin.Context) {
	return []func(ctx *gin.Context){
		AllowCORS,
	}
}
