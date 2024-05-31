package middleware

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/constant"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// a very simple logger, just for demo on how process the request
func Logger(log *logrus.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		statusCode := c.Writer.Status()

		requestId, exist := c.Get(constant.RequestId)
		if !exist {
			requestId = ""
		}

		entry := log.WithFields(logrus.Fields{
			"request_id":  requestId,
			"latency":     time.Since(start),
			"method":      c.Request.Method,
			"status_code": statusCode,
			"path":        path,
		})

		// we responds with server errors
		if statusCode >= 500 && statusCode <= 599 {
			var appErr *apperror.AppError
			for _, err := range c.Errors {
				if errors.As(err, &appErr) {
					entry.WithField("error", appErr).Error("got error")
					entry.WithField("stack", string(appErr.GetStackTrace())).Error("got error")
				}
				// possibly you want to handle other errors
			}

			return
		}

		entry.Info("request processed")
	}

}
