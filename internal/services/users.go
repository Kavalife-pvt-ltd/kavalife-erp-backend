package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/handlers"
	"github.com/paaart/kavalife-erp-backend/internal/models"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
)

func AllUsers(c *gin.Context) {
	data, err := handlers.AllUsers(c)
	if err != nil {
		utils.SuccessWithError(c, err)
		return
	}
	utils.SuccessWithData(c, data)

}

func UserLogin(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, errors.New("invalid JSON input"))
		return
	}
	user, err := handlers.GetUserByUsername(c, req.Username)
	if err != nil {
		if err.Error() == "no rows found" {
			utils.StatusUnauthorized(c, errors.New("invalid username or password"))
		} else {
			utils.StatusInternalServerError(c, err)
		}
		return
	}
	if err := utils.CheckPasswordHash(req.Password, user.Password); !err {
		utils.SuccessWithError(c, errors.New("invalid username or password"))
		return
	}
	token, _ := utils.CreateJWT(user.ID, user.Username, user.Department, user.Role, 24*time.Hour) //creating token
	// c.SetCookie("usrCookie", token, 86400, "/", "", true, true)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "usrCookie",
		Value:    token,
		Path:     "/",
		MaxAge:   86400,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,                  // ✅ Must be true for SameSite=None
		SameSite: http.SameSiteNoneMode, // ✅ Required for cross-site
	})
	user.Password = ""
	utils.SuccessWithData(c, user)
}

func CheckUser(c *gin.Context) {
	userToken, err := c.Cookie("usrCookie")
	if err != nil {
		utils.StatusUnauthorized(c, errors.New("not logged in"))
		return
	}

	claims, valid := utils.ValidateJWT(userToken)
	if !valid {
		utils.StatusUnauthorized(c, errors.New("invalid or expired session"))
		return
	}

	user, err := handlers.GetLoggedUserByUsername(c, claims["username"].(string))
	if err != nil {
		if err.Error() == "user not found" {
			utils.StatusUnauthorized(c, err)
		} else {
			utils.StatusInternalServerError(c, err)
		}
		return
	}
	utils.SuccessWithData(c, user)
}

func Logout(c *gin.Context) {
	// Set the cookie with MaxAge -1 to delete
	c.SetCookie("usrCookie", "", -1, "/", "", true, true)
	utils.SuccessWithData(c, "Logged out successfully")
}
