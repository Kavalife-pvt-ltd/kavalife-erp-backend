package services

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/models"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
)

func AllUsers(c *gin.Context) {
	rows, err := db.DB.Query(c, "select id, username, created_at, role, phone_num FROM public.users")
	if err != nil {
		utils.SuccessWithError(c, err)
		return
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var p models.User
		err := rows.Scan(&p.ID, &p.Username, &p.Created_at, &p.Role, &p.Phone_Num) // Adjust based on your struct
		if err != nil {
			log.Printf("Failed to scan product: %v", err)
			continue
		}
		users = append(users, p)

	}
	if err = rows.Err(); err != nil {
		utils.SuccessWithError(c, err)
		return
	}
	utils.SuccessWithData(c, users)

}

func GetOneUser(c *gin.Context) {
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
	utils.SuccessWithData(c, user)
}

func UserLogin(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, errors.New("invalid JSON input"))
		return
	}
	var user models.User
	err := db.DB.QueryRow(c, `SELECT id, username, role, password, phone_num FROM public.users WHERE username=$1`, req.Username).
		Scan(&user.ID, &user.Username, &user.Role, &user.Password, &user.Phone_Num)

	if err == sql.ErrNoRows {
		utils.StatusUnauthorized(c, errors.New("no data found"))
		return
	} else if err != nil {
		utils.StatusInternalServerError(c, errors.New("query error"))
		return
	}
	if err := utils.CheckPasswordHash(req.Password, user.Password); !err {
		utils.SuccessWithError(c, errors.New("invalid username or password"))
		return
	}
	token, _ := utils.CreateJWT(user.Username, 24*time.Hour) //creating token
	c.SetCookie("usrCookie", token, 86400, "/", "", true, true)
	user.Password = ""
	utils.SuccessWithData(c, user)
}
