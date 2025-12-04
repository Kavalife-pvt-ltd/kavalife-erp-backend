package models

import "time"

// SalesPOStatusLog represents a row in sales_po_status_log.
type SalesPOStatusLog struct {
	ID int64 `db:"id" json:"id"`

	POID int64 `db:"po_id" json:"poId"`

	FromStatus *SalesPOStatus `db:"from_status" json:"fromStatus,omitempty"`
	ToStatus   SalesPOStatus  `db:"to_status" json:"toStatus"`

	ChangedBy *int64 `db:"changed_by" json:"changedBy,omitempty"`

	Note *string `db:"note" json:"note,omitempty"`

	ChangedAt time.Time `db:"changed_at" json:"changedAt"`
}

// CreateStatusLogInput is what your services will pass when recording a transition.
type CreateStatusLogInput struct {
	POID       int64          `json:"poId"`
	FromStatus *SalesPOStatus `json:"fromStatus,omitempty"`
	ToStatus   SalesPOStatus  `json:"toStatus"`
	ChangedBy  *int64         `json:"changedBy,omitempty"`
	Note       *string        `json:"note,omitempty"`
}
