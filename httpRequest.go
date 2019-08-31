package logrusdriver

import (
	"fmt"
	"time"

	"github.com/labstack/echo"
)

// HTTPRequest is HTTPRequest
type HTTPRequest struct {
	RequestMethod string `json:"requestMethod"`
	RequestURL    string `json:"requestUrl"`
	Status        int    `json:"status"`
	UserAgent     string `json:"userAgent"`
	RemoteIP      string `json:"remoteIp"`
	Referer       string `json:"referer"`
	Latency       string `json:"latency"`
	Protocol      string `json:"protocol"`
}

// NewHTTPRequest is NewHTTPRequest
func NewHTTPRequest(c echo.Context, latency time.Duration) *HTTPRequest {
	req := c.Request()
	return &HTTPRequest{
		RequestMethod: req.Method,
		RequestURL:    fmt.Sprintf("%s://%s%s", c.Scheme(), req.Host, req.RequestURI),
		Status:        c.Response().Status,
		UserAgent:     req.UserAgent(),
		RemoteIP:      c.RealIP(),
		Referer:       req.Referer(),
		Latency:       fmt.Sprintf("%fs", latency.Seconds()),
		Protocol:      req.Proto,
	}
}
