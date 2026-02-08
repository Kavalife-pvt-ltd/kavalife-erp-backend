package routes

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/models"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
	"github.com/sirupsen/logrus"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RunApp() *gin.Engine {

	utils.InitLogger() //Initialize logger ONCE at startup

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(GinLoggerMiddleware())
	r.Use(CORSMiddleware())

	r.GET("/", Explorer)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	Routes(r)
	return r
}

var allowedOrigins = map[string]bool{
	"http://localhost:5173":                    true,
	"https://erp.kavalife.in":                  true,
	"https://kavalife-erp-frontend.vercel.app": true,
	"https://kavalife-erp.vercel.app":          true,
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
		}
		c.Next()
	}
}

func GinLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Ignore root and swagger
		if path == "/" || strings.HasPrefix(path, "/swagger") {
			c.Next()
			return
		}

		c.Next()

		latency := time.Since(start)

		entry := utils.Log.WithFields(logrus.Fields{
			"status":    c.Writer.Status(),
			"method":    c.Request.Method,
			"path":      path,
			"latency":   latency.Milliseconds(),
			"client_ip": c.ClientIP(),
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info("request completed")
		}
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := c.Cookie("usrCookie")
		if err != nil {
			utils.StatusUnauthorized(c, errors.New("session expired or not logged in"))
			c.Abort()
			return
		}

		jwtValidate, ok := utils.ValidateJWT(userId)
		if !ok {
			utils.StatusUnauthorized(c, errors.New("user not valid"))
			c.Abort()
			return
		}

		var user models.User
		err = db.DB.QueryRow(
			c,
			`SELECT id, username, role, department, department_role 
			 FROM public.users WHERE id=$1`,
			jwtValidate["id"],
		).Scan(
			&user.ID,
			&user.Username,
			&user.Role,
			&user.Department,
			&user.DepartmentRole,
		)

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
		c.Set("user_id", user.ID)
		c.Set("role", user.Role)
		c.Set("department", user.Department)

		c.Next()
	}
}
