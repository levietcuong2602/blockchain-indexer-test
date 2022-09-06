package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const preflightCacheMaxAge = 12 * time.Hour

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "GET", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token",
			"Authorization", "accept", "origin", "Cache-Control", "X-Requested-With",
		},
		AllowCredentials: true,
		MaxAge:           preflightCacheMaxAge,
	})
}
