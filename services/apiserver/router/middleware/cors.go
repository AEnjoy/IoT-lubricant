package middleware

import (
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/gin-gonic/gin"
)

func AllowCORS(c *gin.Context) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin") //请求头部
	if origin != "" {
		//接收客户端发送的origin （重要！）
		//Access-Control-Allow-Origin是必须的,他的值要么是请求Origin字段的值,要么是一个*, 表示接受任意域名的请求
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		//服务器支持的所有跨域请求的方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE,UPDATE")
		//允许跨域设置可以返回其他子段，可以自定义字段
		//该字段可选。CORS请求时，XMLHttpRequest对象的getResponseHeader()方法只能拿到6个基本字段：Cache-Control、Content-Language、Content-Type、Expires、Last-Modified、Pragma。
		//如果想拿到其他字段，就必须在Access-Control-Expose-Headers里面指定。上面的例子指定，getResponseHeader('FooBar')可以返回FooBar字段的值。
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token, session, token")
		// 允许浏览器（客户端）可以解析的头部 （重要）
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		//设置缓存时间
		//该字段可选，用来指定本次预检请求的有效期，单位为秒。有效期是20天（1728000秒），即允许缓存该条回应1728000秒（即20天），在此期间，不用发出另一条预检请求。
		c.Header("Access-Control-Max-Age", "172800")
		//允许客户端传递校验信息比如 cookie (重要)
		c.Header("Access-Control-Allow-Credentials", "true")
	}
	//允许类型校验
	if method == "OPTIONS" {
		logger.Info("OPTIONS")
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("Panic info is: %v", err)
		}
	}()
	//c.Next()
}
