package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	otpapp "github.com/ppeymann/top-app.git"
)

func (s *Server) redisRateLimit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Skip rate limiting if the path is in the exclude list
		for _, v := range s.Config.RateLimit.RateLimitExcludePaths {
			if strings.HasPrefix(ctx.Request.URL.Path, v) {
				ctx.Next()
				return
			}
		}

		// Get the client IP address
		clientIP := ctx.ClientIP()

		// Generate a unique key for the client's rate limit
		key := fmt.Sprintf("api_rate_limit:%s", clientIP)

		// Check if the client has exceeded the rate limit
		count, err := s.redis.Incr(context.Background(), key).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, otpapp.ErrInternalServer.Error())
			ctx.Abort()
			return
		}

		// Set expiration if it's the first request within the duration
		if count == 1 {
			err := s.redis.Expire(context.Background(), key, time.Duration(s.Config.RateLimit.RateLimitDurationSeconds)*time.Second).Err()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, otpapp.ErrInternalServer.Error())
				ctx.Abort()
				return
			}
		}

		// Check if the request exceeds the allowed limit
		if count > s.Config.RateLimit.RateLimitRequestPerDuration {
			ctx.JSON(http.StatusTooManyRequests, "rate limit exceeded")
			ctx.Abort()
			return
		}

		// Continue processing the request
		ctx.Next()
	}
}
