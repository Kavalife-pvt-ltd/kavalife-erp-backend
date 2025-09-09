// internal/models/vir.go
package models

import "time"

type VIR struct {
	ID          int64             `json:"id"`
	CreatedAt   time.Time         `json:"created_at"`
	VendorID    int64             `json:"vendor_id"`
	ProductID   int64             `json:"product_id"`
	Checklist   map[string]string `json:"checklist"`
	Remarks     string            `json:"remarks"`
	CreatedBy   int64             `json:"created_by"`
	CheckedBy   *int64            `json:"checked_by,omitempty"`
	CheckedAt   *time.Time        `json:"checked_at,omitempty"`
	Status      string            `json:"status"`
	VIRNumber   string            `json:"vir_number"`
	VendorName  *string           `json:"vendor_name,omitempty"`
	ProductName *string           `json:"product_name,omitempty"`
}

type CreateVIRRequest struct {
	Vendor    string            `json:"vendor" binding:"required"`
	Product   string            `json:"product" binding:"required"`
	Checklist map[string]string `json:"checklist" binding:"required"`
	Remarks   string            `json:"remarks"`
	CreatedBy CreatedByData     `json:"createdBy" binding:"required"`
	CreatedAt string            `json:"createdAt" binding:"required"`
}

type CreatedByData struct {
	Data UserData `json:"data" binding:"required"`
}

type UserData struct {
	ID        int64  `json:"id" binding:"required"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	Role      string `json:"role"`
	PhoneNum  *int64 `json:"phone_num,omitempty"`
}

type CreateVIRResponse struct {
	VIRNumber string `json:"vir_number"`
	CreatedBy string `json:"createdBy"`
	CreatedAt string `json:"createdAt"`
	Status    string `json:"status"`
}

type VerifyVIRRequest struct {
	CheckedBy CheckedByData `json:"checkedBy" binding:"required"`
	CheckedAt string        `json:"checkedAt" binding:"required"`
}

type CheckedByData struct {
	Data UserData `json:"data" binding:"required"`
}

type VerifyVIRResponse struct {
	VIRNumber string `json:"vir_number"`
	CheckedBy string `json:"checkedBy"`
	CheckedAt string `json:"checkedAt"`
	Status    string `json:"status"`
}
