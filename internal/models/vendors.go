package models

import "time"

type Vendors struct {
	ID         int       `mapstructure:"id" json:"id"`
	Name       string    `mapstructure:"name" json:"name"`
	Created_at time.Time `mapstructure:"created_at" json:"created_at"`
	Status     string    `mapstructure:"status" json:"status"`
	GovId      string    `mapstructure:"gov_id" json:"gov_id"`
	Type       string    `mapstructure:"type" json:"type"`
}
