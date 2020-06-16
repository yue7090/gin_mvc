package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"github.com/ratelimit"
	"gopkg.in/ini.v1"
	"os"
	"strconv"
)

type Config struct {
	Duration int64
	RateLimit int64
	LimitFunc func(c *gin.Context, ip string)
}

func DefaultLimit() gin.HandlerFunc {
	cfg, err := ini.Load("config/conf.ini")
	if err != nil {
        fmt.Printf("Fail to read file: %v", err)
        os.Exit(1)
	}
	duration := cfg.Section("ratelimit").Key("duration").String()
	rateLimit := cfg.Section("ratelimit").Key("rateLimit").String()
	durationInt, err := strconv.ParseInt(duration,10,64)
	if err != nil {
		fmt.Printf("strconv: %v", err)
		os.Exit(1)
	}
	rateLimitInt, err := strconv.ParseInt(rateLimit,10,64)
	if err != nil {
		fmt.Printf("strconv: %v", err)
		os.Exit(1)
	}
	
	return NewRateLimit(
		Config{
		Duration: durationInt,
		RateLimit: rateLimitInt,
	})
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