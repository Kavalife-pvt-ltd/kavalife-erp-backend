package handlers

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/utils"

	"github.com/paaart/kavalife-erp-backend/internal/models"
)

func RunApp() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(GinLoggerMiddleware()) // default logger
	r.Use(CORSMiddleware())      //CORS might need to update for prod
	Routes(r)
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
		l := utils.InitLogger()
		l.WithFields(map[string]any{
			"status":   c.Writer.Status(),
			"method":   c.Request.Method,
			"path":     path,
			"latency":  latency,
			"clientIP": c.ClientIP(),
		}).Info("Request processed")
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := c.Cookie("usrCookie")
		if err != nil {
			// Cookie is missing or expired
			utils.StatusUnauthorized(c, errors.New("session expired or not logged in"))
			c.Abort()
			return
		}
		jwtValidate, bool := utils.ValidateJWT(userId)
		if !bool {
			utils.StatusUnauthorized(c, errors.New("user not valid"))
			c.Abort()
			return
		}
		var user models.User
		err = db.DB.QueryRow(c, `SELECT id, username, role FROM public.users WHERE username=$1`, jwtValidate["username"]).
			Scan(&user.ID, &user.Username, &user.Role)
		if err == sql.ErrNoRows {
			utils.StatusUnauthorized(c, errors.New("invalid username or password"))
			c.Abort()
			return
		} else if err != nil {
			utils.StatusInternalServerError(c, errors.New("database error"))
			c.Abort()
			return
		}

		c.Set("username", user.Username)
		c.Set("id", user.ID)
		c.Next()
	}
}
