package ratelimit

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok-demo/logger"
)

// connLimiter 并发限流
type connLimiter struct {
	connCurrentConn int
	bucket          chan int
}

func newConnLimiter(cc int) *connLimiter {
	return &connLimiter{
		connCurrentConn: cc,
		bucket:          make(chan int, cc),
	}
}

func (cl *connLimiter) GetConn() bool {
	if len(cl.bucket) >= cl.connCurrentConn {
		logger.Log.Warn("已达到流控上限")
		return false
	}
	cl.bucket <- 1
	return true
}

func (cl *connLimiter) ReleaseConn() {
	<-cl.bucket
	logger.Log.Info("完成任务，释放连接")
}

func RateLimiter(cc int) func(c *gin.Context) {
	connLimiter := newConnLimiter(cc)
	return func(c *gin.Context) {
		defer connLimiter.ReleaseConn()
		if !connLimiter.GetConn() {
			c.JSON(http.StatusOK, "upper limit has been reached")
			c.Abort()
			return
		}
		c.Next()
	}
}
