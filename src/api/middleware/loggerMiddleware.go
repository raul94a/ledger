package middleware

import (
	"fmt"
	"time"
	"src/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)


func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Launching RequestLogger")
		l := app_logger.GetLogger()
		r := c.Request
		// lrw := c.Request.Response
		// bodyBytes, _ := io.ReadAll(c.Request.Body)
		

		// rawBodyString := string(bodyBytes)
		defer func(start time.Time) {
			l.Info(
				fmt.Sprintf(
					"%s request to %s completed",
					r.Method,
					r.RequestURI,
				),
				zap.String("method", r.Method),
				zap.String("url", r.RequestURI),
				zap.String("user_agent", r.UserAgent()),
				// zap.String("body",rawBodyString),
				// zap.Int("status_code", lrw.StatusCode),
				zap.Duration("elapsed_ms", time.Since(start)),
			)
		}(time.Now())

		c.Next()
	}
}
