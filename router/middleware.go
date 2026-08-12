package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/ao9911/go-matrix/log"
	"github.com/ao9911/go-matrix/response"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"

	"github.com/ao9911/bluebell-new/pkg/ecode"
)

const CtxSubjectKey = "userID"

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.JSONFail(c, ecode.NeedLogin, nil)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.JSONFail(c, ecode.InvalidToken, nil)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString
		mc, err := auth.ParseAccess(parts[1])
		if err != nil {
			log.Errorf("jwt.ParseAccessToken error: %v", err)
			response.JSONFail(c, ecode.InvalidToken, nil)
			c.Abort()
			return
		}
		// 将当前请求的 userID 信息保存到请求的上下文 c 上
		c.Set(CtxSubjectKey, mc.Subject)
		c.Next() // 后续的处理函数可以通过 c.Get(CtxSubjectKey) 获取当前请求的用户信息
	}
}

// RateLimitMiddleware 基于令牌桶算法的限流中间件
func RateLimitMiddleware(maxQPS int, maxBurst int) gin.HandlerFunc {
	if maxQPS <= 0 {
		maxQPS = 100
	}
	if maxBurst <= 0 {
		maxBurst = maxQPS
	}

	fillInterval := time.Second / time.Duration(maxQPS)
	bucket := ratelimit.NewBucket(fillInterval, int64(maxBurst))

	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) != 1 {
			c.String(http.StatusTooManyRequests, "rate limit...")
			c.Abort()
			return
		}
		c.Next()
	}
}
