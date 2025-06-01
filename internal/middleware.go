package internal

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/handlers"
	util "github.com/paaart/kavalife-erp-backend/internal/utils"
)

func RunApp() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(GinLoggerMiddleware()) // default logger
	r.Use(CORSMiddleware())      //CORS might need to update for prod
	handlers.Routes(r)
	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func GinLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		l := util.InitLogger()
		l.WithFields(map[string]any{
			"status":   c.Writer.Status(),
			"method":   c.Request.Method,
			"path":     path,
			"latency":  latency,
			"clientIP": c.ClientIP(),
		}).Info("Request processed")
	}
}
