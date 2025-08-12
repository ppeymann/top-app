package server

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PATCH", "OPTION", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authenticate", "Authorization", "X-Requested-With", "Accept", "Accept-Encoding"},
		ExposeHeaders:    []string{"Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
