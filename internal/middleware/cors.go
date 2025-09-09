package middleware

import (
	"time"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-Password-Confirm"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // Must be false when AllowAllOrigins is true
		MaxAge:           12 * time.Hour,
	})
}
