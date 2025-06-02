package models

import "time"

type User struct {
	ID         int       `mapstructure:"id" json:"id"`
	Username   string    `mapstructure:"username" json:"username"`
	Created_at time.Time `mapstructure:"created_at" json:"created_at"`
	Role       string    `mapstructure:"role" json:"role"`
}
