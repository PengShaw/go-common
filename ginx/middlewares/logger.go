package middlewares

import (
	"github.com/PengShaw/go-common/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

const (
	ReqIDContextKey = "req_id"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		id, err := uuid.NewUUID()
		if err != nil {
			logger.Errorf("generate uuid failed: [%s]", err.Error())
			return
		}
		c.Set(ReqIDContextKey, id.String())

		// 处理请求
		c.Next()

		l := logger.WithField("status_code", c.Writer.Status())
		l = l.WithField("latency_time", time.Now().Sub(startTime))
		l = l.WithField("client_ip", c.ClientIP())
		l = l.WithField("req_method", c.Request.Method)
		l = l.WithField("req_uri", c.Request.RequestURI)
		l = l.WithField("req_id", id.String())
		l.Info("")
	}
}
