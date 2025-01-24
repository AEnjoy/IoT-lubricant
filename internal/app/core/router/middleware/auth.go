package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	def "github.com/AEnjoy/IoT-lubricant/pkg/default"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// noneAuth  仅在测试时使用
func noneAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		noneUser := new(casdoorsdk.Claims)
		noneUser.User.Id = uuid.NewString()
		noneUser.User.Name = uuid.NewString()
		c.Set("claims", noneUser)
		c.Next()
	}
}
func casdoorAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		cookie, _ := c.Cookie(model.COOKIE_TOKEY_KEY)
		if cookie != "" {
			token = cookie
		} else {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
				return
			}
			token = parts[1]
		}

		claims, err := casdoorsdk.ParseJwtToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
func Auth() gin.HandlerFunc {
	switch os.Getenv(def.ENV_CORE_AUTH_PROVIDE) {
	case "casdoor":
		return casdoorAuth()
	case "pass", "none", "":
		return noneAuth()
	default:
		panic("unknown auth provider")
	}
}
