package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"github.com/ratelimit"
)

type Config struct {
	Duration int64
	RateLimit int64
	LimitFunc func(c *gin.Context, ip string)
}

var defaultConfig = Config{
	Duration: 60,
	RateLimit:60,
}

func DefaultConfig() Config {
	config := defaultConfig
	return config
}

func DefaultLimit() gin.HandlerFunc {
	return NewRateLimit(defaultConfig)
}

func NewRateLimit(config Config) gin.HandlerFunc {
	ratelimit := ratelimit.NewLimiter(config.Duration, config.RateLimit, time.Second*10)
	if config.LimitFunc == nil {
		config.LimitFunc = func(c *gin.Context, ip string) {
			errorMsg := fmt.Sprintf("rate limit, request should less than %d every %d seconds.", config.RateLimit, config.Duration)
			c.JSON(http.StatusForbidden, gin.H{
				"ip": ip,
				"message": errorMsg,
			})
		}
	}
	return func(c *gin.Context) {
		ip := c.ClientIP()
		shouldLimit := ratelimit.ShouldLimit(ip)
		if shouldLimit {
			config.LimitFunc(c, ip)
			c.Abort()
		}
		c.Next()
	}
}