package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"OPTIONS", "GET", "POST", "PUT"},
		AllowHeaders: []string{"Content-Type", "X-CSRF-Token", "Authorization"},
	})
}
