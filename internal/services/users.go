package services

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/models"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
)

func AllUsers(c *gin.Context) {
	rows, err := db.DB.Query(c, "select * FROM users")
	if err != nil {
		utils.SuccessWithError(c, err)
		return
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Printf("Failed to get values: %v", err)
			continue
		}
		fieldDescriptions := rows.FieldDescriptions()
		rowMap := make(map[string]interface{}, len(values))

		for i, val := range values {
			colName := string(fieldDescriptions[i].Name)
			rowMap[colName] = val
		}

		var user models.User
		if err := mapstructure.Decode(rowMap, &user); err != nil {
			log.Printf("Failed to decode: %v", err)
			continue
		}

		users = append(users, user)
	}
	utils.SuccessWithData(c, users)

}

func UserLogin(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, errors.New("invalid JSON input"))
		return
	}
	var user models.User
	err := db.DB.QueryRow(c, `SELECT id, username, role, password FROM public.users WHERE username=$1`, req.Username).
		Scan(&user.ID, &user.Username, &user.Role, &user.Password)

	if err == sql.ErrNoRows {
		utils.StatusUnauthorized(c, errors.New("invalid username or password"))
		return
	} else if err != nil {
		utils.StatusInternalServerError(c, errors.New("database error"))
		return
	}
	if err := utils.CheckPasswordHash(req.Password, user.Password); !err {
		utils.SuccessWithError(c, errors.New("invalid username or password"))
		return
	}
	token, _ := utils.CreateJWT(user.Username, 24*time.Hour) //creating token
	c.SetCookie("usrCookie", token, 86400, "/", "", true, true)
	utils.SuccessWithData(c, user)
}
