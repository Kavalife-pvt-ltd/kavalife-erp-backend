package models

import "time"

type User struct {
	ID         int       `mapstructure:"id" json:"id"`
	Username   string    `mapstructure:"username" json:"username"`
	Created_at time.Time `mapstructure:"created_at" json:"created_at"`
	Role       string    `mapstructure:"role" json:"role"`
	Password   string    `mapstructure:"password" json:"password"`
	Phone_Num  int       `mapstructure:"phone_num" json:"phone_num"`
}

type LoginRequest struct {
	Username string `mapstructure:"username" json:"username" binding:"required"`
	Password string `mapstructure:"password" json:"password" binding:"required"`
}
