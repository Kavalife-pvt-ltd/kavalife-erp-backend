package models

import (
	"database/sql"
	"time"
)

type Vendors struct {
	ID         int          `mapstructure:"id" json:"id"`
	Name       string       `mapstructure:"name" json:"name"`
	Created_at time.Time    `mapstructure:"created_at" json:"created_at"`
	Status     string       `mapstructure:"status" json:"status"`
	GovId      string       `mapstructure:"gov_id" json:"gov_id"`
	Type       string       `mapstructure:"type" json:"type"`
	UpdatedBy  int          `mapstructure:"updated_by" json:"updated_by,omitempty"`
	Updated_at sql.NullTime `mapstructure:"updated_at" json:"updated_at,omitempty"` //need to check
}
