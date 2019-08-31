package logrusdriver

import (
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

// StackdriverLogging is a echo-middleware suitable for Stackdriver logging using logrus
func StackdriverLogging(config *Config) echo.MiddlewareFunc {
	config = config.withDefault()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			var entry *logrus.Entry
			if config.RequestIDKey == "-" {
				entry = config.Logger.WithFields(nil)
			} else {
				entry = config.Logger.WithField(config.RequestIDKey, config.RequestID(c))
			}

			c = wrapContextWithLogrus(c, entry)
			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			finish := time.Now()

			httpRequest := NewHTTPRequest(c, finish.Sub(start))
			entry.
				WithTime(finish).
				WithFields(logrus.Fields{
					"httpRequest":    httpRequest,
					"requestHeader":  c.Request().Header,
					"responseHeader": c.Response().Header(),
				}).
				Log(config.LogLevel(c), "request completed")
			return nil
		}
	}
}
