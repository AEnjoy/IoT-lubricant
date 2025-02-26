package middleware

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
	def "github.com/AEnjoy/IoT-lubricant/pkg/default"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/services/core/models"
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

var db models.ICoreDb

func casdoorAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		cookie, _ := c.Cookie(model.COOKIE_TOKEY_KEY)
		//logger.Debugf("cookie: %s", cookie)
		if cookie != "" {
			token = cookie
		} else {
			authHeader := c.GetHeader("Authorization")
			logger.Debugf("authHeader: %s", authHeader)
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
			userId, exists := c.Get("userId") // Assuming you have saved userId somewhere in context during login
			if !exists {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				return
			}
			refreshToken, err := db.GetUserRefreshToken(c, userId.(string))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to get refresh token"})
				return
			}
			newToken, err := casdoorsdk.RefreshOAuthToken(refreshToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to refresh token"})
				return
			}

			// Update the cookie with the new access token
			c.SetCookie(model.COOKIE_TOKEY_KEY,
				newToken.AccessToken, int(newToken.Expiry.Unix()-time.Now().Unix()), "/", "",
				false, true)

			// Parse the new token and set claims
			newClaims, err := casdoorsdk.ParseJwtToken(newToken.AccessToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid refreshed token"})
				return
			}
			claims = newClaims
		}

		c.Set("claims", claims)
		c.Set("userId", claims.User.Id)
		c.Next()
	}
}

var _o sync.Once

func Auth() gin.HandlerFunc {
	_o.Do(func() {
		db = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE).(models.ICoreDb)
	})
	switch os.Getenv(def.ENV_CORE_AUTH_PROVIDE) {
	case "casdoor":
		return casdoorAuth()
	case "pass", "none", "":
		return noneAuth()
	default:
		panic("unknown auth provider")
	}
}
