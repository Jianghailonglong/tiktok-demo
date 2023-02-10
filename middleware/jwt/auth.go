package jwt

import (
	"net/http"
	"strings"
	"tiktok-demo/common"

	"github.com/gin-gonic/gin"
)

// AuthInHeader 鉴权中间件，许多接口需要鉴权才能访问
func AuthInHeader() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			responseCodeInvalidMessage(c)
			c.Abort()
			return
		}
		tokenSplits := strings.SplitN(token, " ", 2)
		if len(tokenSplits) != 2 || tokenSplits[0] != "Bearer" {
			responseCodeInvalidMessage(c)
			c.Abort()
			return
		}
		claims, err := parseToken(tokenSplits[1])
		if err != nil {
			responseCodeInvalidMessage(c)
			c.Abort()
			return
		}
		c.Set(ContextUserIDKey, claims.UserId)
		c.Next()
	}
}

// AuthInBody 鉴权中间件
func AuthInBody() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.PostForm("token")
		if token == "" {
			responseCodeInvalidMessage(c)
			c.Abort()
			return
		}
		tokenSplits := strings.SplitN(token, " ", 2)
		if len(tokenSplits) != 2 || tokenSplits[0] != "Bearer" {
			responseCodeInvalidMessage(c)
			c.Abort()
			return
		}
		claims, err := parseToken(tokenSplits[1])
		if err != nil {
			responseCodeInvalidMessage(c)
			c.Abort()
			return
		}
		c.Set(ContextUserIDKey, claims.UserId)
		c.Next()
	}
}

// AuthWithoutLimitLoginStatus 视频流接口需要验证
func AuthWithoutLimitLoginStatus() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Query("token")
		var userId = -1
		if token == "" {
			// token为空，设置userID=-1，继续执行接下来handler
			userId = -1
		} else {
			// token不为空，验证需要通过
			tokenSplits := strings.SplitN(token, " ", 2)
			if len(tokenSplits) != 2 || tokenSplits[0] != "Bearer" {
				userId = -1
				responseCodeInvalidMessage(c)
				c.Abort()
			} else {
				claims, err := parseToken(tokenSplits[1])
				if err != nil {
					userId = -1
					responseCodeInvalidMessage(c)
					c.Abort()
				} else {
					userId = claims.UserId
				}
			}
		}
		c.Set(ContextUserIDKey, userId)
		c.Next()
	}
}

// responseCodeInvalidMessage 提取成一个公共函数
func responseCodeInvalidMessage(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, common.Response{
		StatusCode: -1,
		StatusMsg:  "无效的Token",
	})
}
