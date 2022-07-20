package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s]-%s %s %d %s %s \n",
			params.TimeStamp.Format(time.RFC822),
			params.ClientIP,
			params.Method,
			params.StatusCode,
			params.Latency,
			params.Path,
		)
	})
}
