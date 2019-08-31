package logrusdriver

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

type (
	// RequestID defines request id from context
	RequestID func(c echo.Context) string
	// LogLevel defines log level from context
	LogLevel func(c echo.Context) logrus.Level
)

const defaultRequestIDKey string = "requestID"

// Config is configuration for this logging middleware
type Config struct {
	Logger       *logrus.Logger
	Skipper      middleware.Skipper
	RequestID    RequestID
	RequestIDKey string
	LogLevel     LogLevel
}

func (c *Config) withDefault() *Config {
	if c == nil {
		c = &Config{}
	}
	if c.Logger == nil {
		c.Logger = logrus.New()
	}
	if c.Skipper == nil {
		c.Skipper = middleware.DefaultSkipper
	}
	if c.RequestID == nil {
		c.RequestID = RequestIDFromHeader(echo.HeaderXRequestID)
	}
	if c.RequestIDKey == "" {
		c.RequestIDKey = defaultRequestIDKey
	}
	if c.LogLevel == nil {
		c.LogLevel = LogLevelFromStatus(logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel)
	}
	return c
}

// RequestIDFromHeader defines request id using specified header
func RequestIDFromHeader(header string) RequestID {
	return func(c echo.Context) string {
		id := c.Request().Header.Get(header)
		if id == "" {
			id = c.Response().Header().Get(header)
		}
		return id
	}
}

// LogLevelFromStatus defines log level based on response status
func LogLevelFromStatus(levelFor5xx, levelFor4xx, levelForOther logrus.Level) LogLevel {
	return func(c echo.Context) logrus.Level {
		status := c.Response().Status
		if status >= 500 {
			return levelFor5xx
		}
		if status >= 400 {
			return levelFor4xx
		}
		return levelForOther
	}
}

// LogLevelConstantOf defines same log level for all context
func LogLevelConstantOf(l logrus.Level) LogLevel {
	return func(echo.Context) logrus.Level {
		return l
	}
}
