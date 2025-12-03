package models

import "time"

// Enums as typed strings (nice for autocomplete & avoiding typos)
type SalesPOStatus string
type SalesPORequestType string
type SalesPOFulfillmentType string

const (
	// request_type
	RequestTypeSample   SalesPORequestType = "sample"
	RequestTypePurchase SalesPORequestType = "purchase"

	// status
	StatusQuoteRequested      SalesPOStatus = "quote_requested"
	StatusQuoteAdminApproved  SalesPOStatus = "quote_admin_approved"
	StatusQuoteSentToClient   SalesPOStatus = "quote_sent_to_client"
	StatusClientNegotiation   SalesPOStatus = "client_negotiation"
	StatusClientApproved      SalesPOStatus = "client_approved"
	StatusFinalAdminApproved  SalesPOStatus = "final_admin_approved"
	StatusRoutedToPurchase    SalesPOStatus = "routed_to_purchase"
	StatusRoutedToProduction  SalesPOStatus = "routed_to_production"
	StatusAdminRejected       SalesPOStatus = "admin_rejected"
	StatusClientRejected      SalesPOStatus = "client_rejected"
	StatusCancelled           SalesPOStatus = "cancelled"
	StatusPurchaseCompleted   SalesPOStatus = "purchase_completed"
	StatusProductionCompleted SalesPOStatus = "production_completed"
	StatusClosed              SalesPOStatus = "closed"

	// fulfillment_type
	FulfillmentPurchase   SalesPOFulfillmentType = "purchase"
	FulfillmentProduction SalesPOFulfillmentType = "production"
)

// SalesPO represents a row in the sales_po table.
type SalesPO struct {
	ID int64 `db:"id" json:"id"`

	PONumber  *string `db:"po_number" json:"poNumber,omitempty"`
	ProductID int64   `db:"product_id" json:"productId"`

	CompanyName          string  `db:"company_name" json:"companyName"`
	CompanyAddress       string  `db:"company_address" json:"companyAddress"`
	COAURL               *string `db:"coa_url" json:"coaUrl,omitempty"`
	CompanyContactName   *string `db:"company_contact_name" json:"companyContactName,omitempty"`
	CompanyContactNumber *string `db:"company_contact_number" json:"companyContactNumber,omitempty"`
	CompanyContactEmail  *string `db:"company_contact_email" json:"companyContactEmail,omitempty"`

	Purity *string `db:"purity" json:"purity,omitempty"`
	Grade  *string `db:"grade" json:"grade,omitempty"`

	RequestType  SalesPORequestType `db:"request_type" json:"requestType"`
	Quantity     float64            `db:"quantity" json:"quantity"`
	QuantityUnit *string            `db:"quantity_unit" json:"quantityUnit,omitempty"`
	AskingPrice  *float64           `db:"asking_price" json:"askingPrice,omitempty"`

	Comments             *string    `db:"comments" json:"comments,omitempty"`
	ExpectedDeliveryDate *time.Time `db:"expected_delivery_date" json:"expectedDeliveryDate,omitempty"`
	RequestDate          time.Time  `db:"request_date" json:"requestDate"`

	Status SalesPOStatus `db:"status" json:"status"`

	SendTo string `json:"sendTo,omitempty"`

	SalesRepID      *int64     `db:"sales_rep_id" json:"salesRepId,omitempty"`
	ApprovedByID    *int64     `db:"approved_by_id" json:"approvedById,omitempty"`
	ApprovedAt      *time.Time `db:"approved_at" json:"approvedAt,omitempty"`
	RejectedByID    *int64     `db:"rejected_by_id" json:"rejectedById,omitempty"`
	RejectionReason *string    `db:"rejection_reason" json:"rejectionReason,omitempty"`

	FulfillmentType   *SalesPOFulfillmentType `db:"fulfillment_type" json:"fulfillmentType,omitempty"`
	PurchaseOrderID   *int64                  `db:"purchase_order_id" json:"purchaseOrderId,omitempty"`
	ProductionBatchID *int64                  `db:"production_batch_id" json:"productionBatchId,omitempty"`

	PackedByID   *int64     `db:"packed_by_id" json:"packedById,omitempty"`
	PackedAt     *time.Time `db:"packed_at" json:"packedAt,omitempty"`
	DeliveryCode *string    `db:"delivery_code" json:"deliveryCode,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// CreateSalesPORequest is what the frontend sends when creating a new quote request.
type CreateSalesPORequest struct {
	ProductID int64 `json:"productId"`

	CompanyName          string  `json:"companyName"`
	CompanyAddress       string  `json:"companyAddress"`
	COAURL               *string `json:"coaUrl,omitempty"`
	CompanyContactName   *string `json:"companyContactName,omitempty"`
	CompanyContactNumber *string `json:"companyContactNumber,omitempty"`
	CompanyContactEmail  *string `json:"companyContactEmail,omitempty"`

	Purity *string `json:"purity,omitempty"`
	Grade  *string `json:"grade,omitempty"`

	RequestType  SalesPORequestType `json:"requestType"` // "sample" | "purchase"
	Quantity     float64            `json:"quantity"`
	QuantityUnit *string            `json:"quantityUnit,omitempty"`
	AskingPrice  *float64           `json:"askingPrice,omitempty"`

	Comments             *string    `json:"comments,omitempty"`
	ExpectedDeliveryDate *time.Time `json:"expectedDeliveryDate,omitempty"`
	RequestDate          *time.Time `json:"requestDate,omitempty"` // optional override

	// SalesRepID will usually come from auth context, not from the body,
	// but we keep it here if you ever want to override in services.
	SalesRepID *int64 `json:"salesRepId,omitempty"`
}

// UpdateSalesPOStatusRequest is for status transitions (admin/client actions).
type UpdateSalesPOStatusRequest struct {
	POID            int64                   `json:"poId"`
	ToStatus        SalesPOStatus           `json:"toStatus"`
	NewQuantity     *float64                `json:"newQuantity,omitempty"`
	NewAskingPrice  *float64                `json:"newAskingPrice,omitempty"`
	NewComments     *string                 `json:"newComments,omitempty"`
	RejectionReason *string                 `json:"rejectionReason,omitempty"`
	FulfillmentType *SalesPOFulfillmentType `json:"fulfillmentType,omitempty"`
	DeliveryCode    *string                 `json:"deliveryCode,omitempty"`
	SendTo          *string                 `json:"sendTo,omitempty"` // ✅ NEW
}

// SalesPOResponse is what you might send back to the frontend.
// For now it's just an alias of SalesPO; you can customize later if needed.
type SalesPOResponse = SalesPO
