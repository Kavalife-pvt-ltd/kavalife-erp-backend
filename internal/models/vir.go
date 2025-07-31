package models

import (
	"time"
)

type VIR struct {
	ID        int64                  `mapstructure:"id" json:"id"`
	CreatedAt time.Time              `mapstructure:"created_at" json:"created_at"`
	VendorID  int64                  `mapstructure:"vendor_id" json:"vendor_id"`
	ProductID int64                  `mapstructure:"product_id" json:"product_id"`
	Checklist map[string]interface{} `mapstructure:"checklist" json:"checklist"`
	Remarks   string                 `mapstructure:"remarks" json:"remarks"`
	CreatedBy int64                  `mapstructure:"created_by" json:"created_by"`
	CheckedBy *int64                 `mapstructure:"checked_by" json:"checked_by,omitempty"`
	CheckedAt *time.Time             `mapstructure:"checked_at" json:"checked_at,omitempty"`
	Status    string                 `mapstructure:"status" json:"status"`
	VIRNumber string                 `mapstructure:"vir_number" json:"vir_number"`
}

type CreateVIRRequest struct {
	Vendor    string                 `mapstructure:"vendor" json:"vendor" binding:"required"`
	Product   string                 `mapstructure:"product" json:"product" binding:"required"`
	Checklist map[string]interface{} `mapstructure:"checklist" json:"checklist" binding:"required"`
	Remarks   string                 `mapstructure:"remarks" json:"remarks"`
	CreatedBy CreatedByData          `mapstructure:"createdBy" json:"createdBy"`
	CreatedAt string                 `mapstructure:"createdAt" json:"createdAt"`
}

type CreatedByData struct {
	Data UserData `mapstructure:"data" json:"data" binding:"required"`
}

type UserData struct {
	ID        int64  `mapstructure:"id" json:"id" binding:"required"`
	Username  string `mapstructure:"username" json:"username"`
	CreatedAt string `mapstructure:"created_at" json:"created_at"`
	Role      string `mapstructure:"role" json:"role"`
	PhoneNum  int64  `mapstructure:"phone_num" json:"phone_num"`
}

type VIRResponse struct {
	ID        int64                  `mapstructure:"id" json:"id"`
	VIRNumber string                 `mapstructure:"vir_number" json:"vir_number"`
	CreatedAt time.Time              `mapstructure:"" json:"created_at"`
	VendorID  int64                  `mapstructure:"" json:"vendor_id"`
	ProductID int64                  `mapstructure:"" json:"product_id"`
	Vendor    string                 `mapstructure:"" json:"vendor"`
	Product   string                 `mapstructure:"" json:"product"`
	Checklist map[string]interface{} `mapstructure:"" json:"checklist"`
	Remarks   string                 `mapstructure:"" json:"remarks"`
	CreatedBy int64                  `mapstructure:"" json:"created_by"`
	CheckedBy *int64                 `mapstructure:"" json:"checked_by,omitempty"`
	CheckedAt *time.Time             `mapstructure:"" json:"checked_at,omitempty"`
	Status    string                 `mapstructure:"" json:"status"`
}

type VerifyVIRRequest struct {
	CheckedBy CheckedByData `json:"checkedBy" binding:"required"`
	CheckedAt string        `json:"checkedAt" binding:"required"`
}

type CheckedByData struct {
	Data UserData `json:"data" binding:"required"`
}

type CreateVIRResponse struct {
	VIRNumber string `json:"vir_number"`
	CreatedBy string `json:"createdBy"`
	CreatedAt string `json:"createdAt"` // Use string for ISO format
	Status    string `json:"status"`
}

type VerifyVIRResponse struct {
	VIRNumber string `json:"vir_number"`
	CheckedBy string `json:"checkedBy"`
	CheckedAt string `json:"checkedAt"` // Use string for ISO format
	Status    string `json:"status"`
}
