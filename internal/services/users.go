package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/handlers"
	"github.com/paaart/kavalife-erp-backend/internal/models"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
	"golang.org/x/sync/errgroup"
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

func CreateNewUser(c *gin.Context) {
	var req models.NewUser
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, errors.New("invalid JSON input"))
		return
	}
	// check number
	if len(req.Mobile) != 10 {
		utils.BadRequest(c, errors.New("mobile number must be exactly 10 digits"))
		return
	}
	if len(req.Name) == 0 {
		utils.BadRequest(c, errors.New("Please provide name"))
		return
	}
	//convert password
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.SuccessWithError(c, err)
		return
	}
	req.Password = hashPassword
	//insert in table
	id, err := handlers.InsertNewUser(c, req)
	if err != nil {
		utils.SuccessWithError(c, err)
		return
	}
	utils.SuccessWithData(c, fmt.Sprintf("User inserted with id %d", id))
}

func AllNewUsersList(c *gin.Context) {
	data, err := handlers.AllNewUser(c)
	if err != nil {
		utils.SuccessWithError(c, err)
		return
	}
	utils.SuccessWithData(c, data)
}
func checkUserExists(ctx context.Context, username, phone string) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		user, err := handlers.GetUserByUsername(ctx, username)
		if err != nil {
			if err.Error() == "no rows found" {
				return nil
			}
			return err
		}
		if user != nil {
			return errors.New("username already exists")
		}
		return nil
	})

	g.Go(func() error {
		user, err := handlers.GetUserByPhoneNum(ctx, phone)
		if err != nil {
			if err.Error() == "no rows found" {
				return nil
			}
			return err
		}
		if user != nil {
			return errors.New("phone number already exists")
		}
		return nil
	})

	return g.Wait()
}

func ApproveNewUser(c *gin.Context) {
	var req models.ApproveNewUsers
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, errors.New("invalid JSON input"))
		return
	}
	userID, exists := c.Get("id")
	if !exists {
		utils.BadRequest(c, errors.New("user ID not found in context"))
		return
	}
	uid, ok := userID.(int)
	if !ok {
		utils.BadRequest(c, errors.New("invalid user ID format"))
		return
	}
	role, exists := c.Get("role")
	if !exists {
		utils.BadRequest(c, errors.New("user role not found in context"))
		return
	}
	urole, ok := role.(string)
	if !ok {
		utils.BadRequest(c, errors.New("invalid user role format"))
		return
	}
	if urole != "admin" {
		utils.SuccessWithError(c, errors.New("User not allowed to do the action"))
		return
	}
	//find newUser via id
	newUser, err := handlers.FindNewUserByID(c, req.ID)
	if err != nil {
		utils.SuccessWithError(c, err)
		return
	}
	// log.Println("------", newUser)
	//find username and number in users if found then return with already exists
	if err := checkUserExists(c, newUser.Username, newUser.Mobile); err != nil {
		utils.SuccessWithError(c, err)
		return
	}

	if err != nil && err.Error() != "no rows in result set" {
		utils.SuccessWithError(c, err)
		return
	}
	log.Println("------", newUser)
	//insert new Users and delete from newUsers
	newId, err := handlers.CreateUsers(c, *newUser, req.Role, req.Department, uid)
	if err != nil {
		utils.SuccessWithError(c, err)
		return
	}
	//delete NewUser
	err = handlers.DeleteNewUser(c, req.ID)
	if err != nil {
		utils.SuccessWithError(c, err)
		return
	}

	utils.SuccessWithData(c, fmt.Sprintf("User inserted with id %d", newId))

}
func Logout(c *gin.Context) {
	// Set the cookie with MaxAge -1 to delete
	c.SetCookie("usrCookie", "", -1, "/", "", true, true)
	utils.SuccessWithData(c, "Logged out successfully")
}
