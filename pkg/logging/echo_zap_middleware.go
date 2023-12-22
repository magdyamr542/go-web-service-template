package logging

import (
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerMigglewareConfig struct {
	// Paths to skip
	Skip []string
}

// ZapLogger is an example of echo middleware that logs requests using logger "zap"
func ZapLogger(cfg LoggerMigglewareConfig, log *zap.Logger) echo.MiddlewareFunc {
	log = log.Named("echo_http")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if shouldSkip(cfg.Skip, c.Request().RequestURI) {
				log.Debug("Skip logging request", zap.String("request", c.Request().Method+" "+c.Request().RequestURI))
				err := next(c)
				if err != nil {
					c.Error(err)
				}
				return nil
			}

			req := c.Request()
			res := c.Response()

			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			fields := []zapcore.Field{
				zap.String("remote_ip", c.RealIP()),
				zap.String("latency", time.Since(start).String()),
				zap.String("host", req.Host),
				zap.String("request", req.Method+" "+req.RequestURI),
				zap.Int("status", res.Status),
				zap.Int64("size", res.Size),
				zap.String("user_agent", req.UserAgent()),
				zap.String("request_id", id),
			}

			n := res.Status
			switch {
			case n >= 500:
				log.With(zap.Error(err)).Error("Server error", fields...)
			case n >= 400:
				log.With(zap.Error(err)).Warn("Client error", fields...)
			case n >= 300:
				log.Info("Redirection", fields...)
			default:
				log.Info("Success", fields...)
			}
			return nil
		}
	}
}

func shouldSkip(skip []string, value string) bool {
	for _, prefix := range skip {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}
