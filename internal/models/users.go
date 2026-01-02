package models

import "time"

type User struct {
	ID             int       `mapstructure:"id" json:"id"`
	Username       string    `mapstructure:"username" json:"username"`
	Created_at     time.Time `mapstructure:"created_at" json:"created_at"`
	Name           string    `mapstructure:"name" json:"name"`
	Email          string    `mapstructure:"email" json:"email"`
	Role           string    `mapstructure:"role" json:"role"`
	Password       string    `mapstructure:"password" json:"password,omitempty"`
	Phone_Num      string    `mapstructure:"phone_num" json:"phone_num"`
	Department     string    `json:"department"`
	DepartmentRole *string   `json:"departmentRole,omitempty"`
}

type NewUser struct {
	ID       int    `mapstructure:"id" json:"id"`
	Username string `mapstructure:"username" json:"username" binding:"required"`
	Name     string `mapstructure:"name" json:"name" binding:"required"`
	Email    string `mapstructure:"email" json:"email"`
	Password string `mapstructure:"password" json:"password" binding:"required"`
	Mobile   string `mapstructure:"mob_number" json:"mob_number" binding:"required"`
}
type LoginRequest struct {
	Username string `mapstructure:"username" json:"username" binding:"required"`
	Password string `mapstructure:"password" json:"password" binding:"required"`
}

type LoggedUserDetails struct {
	ID             int     `json:"id"`
	Username       string  `json:"username"`
	Role           string  `json:"role"`
	Department     string  `json:"department"`
	DepartmentRole *string `json:"departmentRole,omitempty"`
}
type ApproveNewUsers struct {
	ID         int    `mapstructure:"id" json:"id"`
	Role       string `mapstructure:"role" json:"role"`
	Department string `mapstructure:"department" json:"department"`
	// DepartmentRole *string `json:"departmentRole,omitempty"`
}

func (v User) GetID() int {
	return v.ID
}

func (v NewUser) GetID() int {
	return v.ID
}
